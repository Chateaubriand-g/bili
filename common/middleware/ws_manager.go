package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var wsupgrader = websocket.Upgrader{}

type WSClient struct {
	Conn     *websocket.Conn
	UserID   uint64
	Send     chan []byte
	LastPong time.Time
}

type WSManager struct {
	Mtx     sync.RWMutex
	Clients map[uint64]map[*WSClient]bool //userID:all client of this id
	Ctx     context.Context
	Cancel  context.CancelFunc
	rds     *redis.Client
}

// *WSManager
func NewWSManager(ctx context.Context, rds *redis.Client) *WSManager {
	res := &WSManager{
		Clients: make(map[uint64]map[*WSClient]bool),
		rds:     rds,
	}
	res.Ctx, res.Cancel = context.WithCancel(ctx)
	return res
}

func (m *WSManager) AddClient(c *WSClient) {
	m.Mtx.Lock()
	defer m.Mtx.Unlock()

	clients, exists := m.Clients[c.UserID]
	if !exists {
		clients = make(map[*WSClient]bool)
		m.Clients[c.UserID] = clients
	}
	if len(clients) > 5 {
		log.Println("user:", c.UserID, "attempt to has more than 5 wsclient")
		return
	}
	clients[c] = true
}

func (m *WSManager) RemoveClient(c *WSClient) {
	m.Mtx.Lock()
	defer m.Mtx.Unlock()

	if clients, exists := m.Clients[c.UserID]; exists {
		delete(clients, c)
		c.Conn.Close()
		close(c.Send)
		if len(clients) == 0 {
			delete(m.Clients, c.UserID)
		}
	}
}

func (m *WSManager) SendToUser(userID uint64, msg []byte) {
	m.Mtx.Lock()
	defer m.Mtx.Unlock()

	if clients, exists := m.Clients[userID]; exists {
		for c := range clients {
			select {
			case c.Send <- msg:
			default:
				go m.RemoveClient(c)
			}
		}
	}
}

// 订阅redis的"notifu:user:*"pubsub，
// 当rocketmqconsumer推送消息给redis时，redis会给pubsub对应的channel发送消息，订阅者可以通过该channel获取
func (m *WSManager) SubscribeRedis() {
	pubsub := m.rds.PSubscribe(context.Background(), "notify:user:*")
	ch := pubsub.Channel()

	for {
		select {
		case <-m.Ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			var uid uint64
			fmt.Sscanf(msg.Channel, "notify:user:%d", &uid)
			m.SendToUser(uid, []byte(msg.Payload))
		}
	}
}

func (m *WSManager) ReadPump(c *WSClient) {
	defer m.RemoveClient(c)

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.LastPong = time.Now()
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (m *WSManager) WritePump(c *WSClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		m.RemoveClient(c)
	}()

	for {
		select {
		case <-m.Ctx.Done():
			return
		//从send通道接收到redis的待推送消息
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		//定时器触发
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func WSHandler(w http.ResponseWriter, r *http.Request, userID uint64, wsmanager *WSManager) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrader err:", err)
		return
	}

	client := WSClient{
		Conn:     conn,
		UserID:   userID,
		Send:     make(chan []byte, 256),
		LastPong: time.Now(),
	}

	wsmanager.AddClient(&client)
}
