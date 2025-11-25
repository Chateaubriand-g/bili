package proxy

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	grpcConnCache   sync.Map
	cacheExpireTime = 30 * time.Minute
	cacheMaxNums    = 100
	connCount       int32
)

func GetGRPCConn(cli *api.Client, serverName string) (*grpc.ClientConn, error) {
	server, err := PickService(cli, serverName)
	if err != nil {
		return nil, fmt.Errorf("pickservice error: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", server.Address, server.Port)

	if conn, exists := grpcConnCache.Load(addr); exists {
		c := conn.(*grpc.ClientConn)
		if isConnAvailable(c) {
			return c, nil
		}
		grpcConnCache.Delete(conn)
		c.Close()
		atomic.AddInt32(&connCount, -1)
	}

	if connCount >= int32(cacheMaxNums) {
		return nil, errors.New("grpc conn is full")
	}

	conn, err := createGRPCConn(addr)
	if err != nil {
		return nil, fmt.Errorf("create grpc conn err: %w", err)
	}

	if target, loaded := grpcConnCache.LoadOrStore(addr, conn); loaded {
		conn.Close()
		return target.(*grpc.ClientConn), nil
	}

	atomic.AddInt32(&connCount, 1)
	return conn, nil
}

func createGRPCConn(target string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceConfig := `{
		// "loadBalancingConfig": [{"round_robin": {}}],  // 负载均衡（可选）
		"keepaliveConfig": {
			"keepaliveTime": 10,        // 保活探测间隔（秒）
			"keepaliveTimeout": 3,      // 探测超时时间（秒）
			"permitWithoutStream": true, // 无数据流时允许保活
			"maxConnectionIdle": 900,   // 连接最大空闲时间(秒,15分钟)
			"maxConnectionAge": 1800    // 连接最大存活时间(秒,30分钟)
		}
	}`

	// 保活配置（自动重连+空闲超时）
	kaParams := keepalive.ClientParameters{
		Time:                10 * time.Second, // 每10秒发送一次保活探测
		Timeout:             3 * time.Second,  // 探测超时时间
		PermitWithoutStream: true,             // 无数据流时也发送保活（避免空闲连接被断开）
	}

	// 3. 建立连接（带阻塞等待就绪+超时）
	cc, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithKeepaliveParams(kaParams), // 保活+自动重连
		grpc.WithBlock(),                   // 阻塞直到连接就绪（结合超时）
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*10), // 最大接收消息大小（10MB，根据业务调整）
		),
	)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func isConnAvailable(conn *grpc.ClientConn) bool {
	switch conn.GetState().String() {
	case "READY", "IDLE":
		return true
	default:
		conn.ResetConnectBackoff()
		return false
	}
}
