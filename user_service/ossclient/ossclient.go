package ossclient

import (
	"time"

	"github.com/Chateaubriand-g/bili/user_service/config"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func InitOssClient(cfg *config.Config) *oss.Client {
	osscfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion("cn-guangzhou").
		WithConnectTimeout(10 * time.Second).
		WithReadWriteTimeout(30 * time.Second)

	client := oss.NewClient(osscfg)
	return client
}
