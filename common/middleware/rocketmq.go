package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/Chateaubriand-g/bili/common/model"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/redis/go-redis/v9"
	"go.uber.org/ratelimit"
	"gorm.io/gorm"
)

type MQMsg struct {
	UserID     uint64                 `json:"userid"`
	Type       int                    `json:"type"`
	FromUserID uint64                 `json:"from_userid,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
}

var msgLimiter = ratelimit.New(2000)

func RegisterConsumer(cfg *config.Config, db *gorm.DB, rds *redis.Client) error {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{cfg.RocketMQ.NameServer}),
		consumer.WithGroupName("notification-consumer-group"),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithMaxReconsumeTimes(3),
		consumer.WithConsumeGoroutineNums(20),
	)
	if err != nil {
		return fmt.Errorf("create consumer failed: %w", err)
	}

	if err := c.Subscribe(
		"notifications",
		consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				processCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				msgLimiter.Take()

				if err := processSingleMessage(processCtx, msg, db, rds); err != nil {
					log.Printf("message consume failed,try again: msgID=%s,err=%v", msg.MsgId, err)
					return consumer.ConsumeRetryLater, err
				}
			}
			return consumer.ConsumeSuccess, nil
		},
	); err != nil {
		return fmt.Errorf("subscribe failed: %w", err)
	}

	if err := c.Start(); err != nil {
		return fmt.Errorf("consumer start failed: %w", err)
	}

	consumerShutdownHook(c)
	return nil
}

func processSingleMessage(ctx context.Context, msg *primitive.MessageExt, db *gorm.DB, rds *redis.Client) error {
	var m MQMsg
	if err := json.Unmarshal(msg.Body, &m); err != nil {
		log.Printf("conumer msg unmarshal failed")
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	if m.UserID == 0 || m.Type < 0 || m.Type > 5 {
		return fmt.Errorf("invaild msg:userID = %d,type=%d", m.UserID, m.Type)
	}

	newmsg, err := saveToDB(ctx, &m, db)
	if err != nil {
		return fmt.Errorf("msg save to db failed: %w", err)
	}

	if err := updateRedisAndPublish(ctx, newmsg, &m, rds); err != nil {
		return fmt.Errorf("update redis failed: %w", err)
	}
	return nil
}

func saveToDB(ctx context.Context, m *MQMsg, db *gorm.DB) (*model.Notification, error) {
	newmsg := model.Notification{
		UserID:     m.UserID,
		Type:       int8(m.Type),
		IsRead:     0,
		FromUserID: m.FromUserID,
	}

	if m.Payload != nil {
		payloadBytes, err := json.Marshal(m.Payload)
		if err != nil {
			return nil, fmt.Errorf("payload marshal failed: %w", err)
		}
		newmsg.Payload = string(payloadBytes)
	} else {
		newmsg.Payload = "{}"
	}

	var err error
	for i := 0; i < 3; i++ {
		if err := db.WithContext(ctx).Create(&newmsg).Error; err == nil {
			return &newmsg, nil
		}
		time.Sleep(time.Millisecond * 100 * (1 << i))
	}
	return nil, fmt.Errorf("msg save to db failed: %w", err)
}

func updateRedisAndPublish(ctx context.Context, msg *model.Notification, mqmsg *MQMsg, rds *redis.Client) error {
	totalkey := fmt.Sprintf("notify:unread:%d", mqmsg.UserID)
	typekey := fmt.Sprintf("notify:unread:type:%d:%d", mqmsg.UserID, mqmsg.Type)
	ch := fmt.Sprintf("notify:user:%d", mqmsg.UserID)

	pushData := map[string]interface{}{
		"id":         msg.ID,
		"user_id":    msg.UserID,
		"type":       msg.Type,
		"payload":    msg.Payload,
		"created_at": msg.CreatedAt.Unix(),
	}

	pushBytes, err := json.Marshal(pushData)
	if err != nil {
		return fmt.Errorf("updateRedisAndPublish==>pushdata marshal failed: %w", err)
	}

	pipe := rds.Pipeline()
	pipe.Incr(ctx, totalkey)
	pipe.Incr(ctx, typekey)
	pipe.Publish(ctx, ch, pushBytes)
	//事务执行上面的命令，要么全部成功，要么全部失败
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("updateRedisAndPublish==>pipe exec failed: %w", err)
	}
	return nil
}

func consumerShutdownHook(c rocketmq.PushConsumer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan

		c.Shutdown()
	}()
}
