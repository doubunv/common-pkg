package v1

import (
	"context"
	"encoding/json"
	"github.com/doubunv/common-pkg/amqp/kafka/config"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logc"
	"runtime/debug"
)

type MessageHandle func(ctx context.Context, msg string) error // 消息处理回调handle

// Consumer   消费者结构体
type Consumer struct {
	reader      *kafka.Reader
	isUserClose bool // 是否是用户主动关闭
}

// NewConsumer   生成一个新的消费者
func NewConsumer(conf config.CustomerConfig) *Consumer {
	kafkaConf := kafka.ReaderConfig{
		Brokers:        conf.Brokers,
		GroupID:        conf.GroupID,
		Topic:          conf.Topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 0,    // 禁用自动提交
		StartOffset:    kafka.FirstOffset,
	}
	reader := kafka.NewReader(kafkaConf)
	return &Consumer{
		reader:      reader,
		isUserClose: false,
	}
}

// AckMessage      消息确认
func (c *Consumer) AckMessage(msg kafka.Message) error {
	return c.reader.CommitMessages(context.Background(), msg)
}

// Close           关闭kafka.Reader,同时退出消费协程
func (c *Consumer) Close() error {
	c.isUserClose = true
	return c.reader.Close()
}

func (c *Consumer) ConsumeMessagesWithContext(handler MessageHandle) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logc.Errorf(context.Background(), "ConsumeMessagesWithContext error:%v, %s", err, string(debug.Stack()))
			}
		}()
		for {
			msg, err := c.reader.ReadMessage(context.Background())
			newCtx := context.Background()
			logc.Infof(newCtx, "---- kafka:ConsumeMessagesWithContext:topic: %s, msg: %s", msg.Topic, string(msg.Value))

			ka := &KafkaMessage{}
			err = json.Unmarshal(msg.Value, ka)
			if err != nil {
				logc.Errorf(newCtx, "---- kafka:ConsumeMessagesWithContext:topic: %s, err: %+v", msg.Topic, err)
				c.AckMessage(msg)
				continue
			}
			newCtx = ka.SetContext(newCtx)
			//go func(msg kafka.Message) {
			//	defer func() {
			//		if err := recover(); err != nil {
			//			logc.Errorf(context.Background(), "ConsumeMessagesWithContext handler error:%v, %s, %s", string(msg.Value), err, string(debug.Stack()))
			//			c.AckMessage(msg)
			//		}
			//	}()
			//	err = handler(newCtx, ka.GetMsg())
			//	if err == nil {
			//		c.AckMessage(msg)
			//	}
			//}(msg)
			err = handler(newCtx, ka.GetMsg())
			if err == nil {
				c.AckMessage(msg)
			} else {
				logc.Errorf(newCtx, "---- kafka:ConsumeMessagesWithContext:topic: %s, err: %+v", msg.Topic, err)
			}
			continue
		}
	}()
	select {}
}
