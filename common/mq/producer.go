package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Chateaubriand-g/bili/common/config"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type RocketMQProducer struct {
	producer rocketmq.Producer
	topic    string
	mtx      sync.RWMutex
}

func NewProducer(cfg *config.Config, topic string) (*RocketMQProducer, error) {
	producer, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{cfg.RocketMQ.NameServer}),
		producer.WithRetry(2),
	)
	if err != nil {
		return nil, fmt.Errorf("produer create err: %w", err)
	}

	if err := producer.Start(); err != nil {
		return nil, fmt.Errorf("producer start err: %w", err)
	}

	return &RocketMQProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *RocketMQProducer) SendEvent(ctx context.Context, v interface{}) error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal v err: %w", err)
	}

	msg := &primitive.Message{
		Topic: p.topic,
		Body:  body,
	}

	_, err = p.producer.SendSync(ctx, msg)
	if err != nil {
		return fmt.Errorf("sendsync err: %w", err)
	}
	return nil
}

func (p *RocketMQProducer) Shutdown() error {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.producer != nil {
		err := p.producer.Shutdown()
		if err != nil {
			return err
		}
	}

	p.producer = nil
	return nil
}
