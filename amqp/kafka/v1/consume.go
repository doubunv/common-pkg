package v1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/doubunv/common-pkg/amqp/kafka/config"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logc"
	"runtime/debug"
	"strings"
	"time"
)

type MessageHandle func(ctx context.Context, msg string) error // 消息处理回调handle

// Consumer   消费者结构体
type Consumer struct {
	reader      *kafka.Reader
	isUserClose bool // 是否是用户主动关闭
	kafkaConf   kafka.ReaderConfig
}

// NewConsumer   生成一个新的消费者
func NewConsumer(conf config.CustomerConfig) *Consumer {
	kafkaConf := kafka.ReaderConfig{
		Brokers:     conf.Brokers,
		GroupID:     conf.GroupID,
		Topic:       conf.Topic,
		StartOffset: conf.StartOffset,
	}
	reader := kafka.NewReader(kafkaConf)
	return &Consumer{
		reader:      reader,
		isUserClose: false,
		kafkaConf:   kafkaConf,
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

func (c *Consumer) sendDeadLetterQueue(ctx context.Context, topic string, msg *KafkaMessage) {
	headT := "mq_dead_letter:"
	if !strings.HasPrefix(topic, headT) {
		topic = headT + topic
	}

	cf := config.ProviderConfig{
		Brokers: c.kafkaConf.Brokers,
		Topic:   topic,
	}
	NewProducer(cf).ProduceMessageWithContext(ctx, msg)
}

func (c *Consumer) handleMessage(ctx context.Context, handler MessageHandle, msg kafka.Message) {
	defer func() {
		if r := recover(); r != nil {
			logc.Errorf(ctx, "Panic recovered in handler: %v\n%s", r, string(debug.Stack()))
		}
	}()

	if len(msg.Value) == 0 {
		return
	}

	ka := &KafkaMessage{}
	if err := json.Unmarshal(msg.Value, ka); err != nil {
		logc.Errorf(ctx, "Unmarshal message error: %v, msg: %s", err, string(msg.Value))
		return
	}

	msgCtx := ka.SetContext(ctx)

	const maxRetries = 3
	delay := time.Second
	for i := 1; i <= maxRetries; i++ {
		if err := handler(msgCtx, ka.GetMsg()); err != nil {
			logc.Infof(msgCtx, "Handle message error (try %d/%d): %v", i, maxRetries, err)
			time.Sleep(delay)
			delay *= 2
			continue
		}
		return
	}

	logc.Errorf(msgCtx, "Message failed after retries, send to DLQ: %s", string(msg.Value))
	// c.sendDeadLetterQueue(msgCtx, msg.Topic, ka)
}

func (c *Consumer) ConsumeMessagesWithContext(ctx context.Context, handler MessageHandle) {
	defer logc.Error(ctx, "MQ consumer stopped")
	const workerCount = 100
	msgCh := make(chan kafka.Message, 100)

	for i := 0; i < workerCount; i++ {
		go func() {
			for msg := range msgCh {
				c.handleMessage(ctx, handler, msg)
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			close(msgCh)
			logc.Info(ctx, "Context canceled, exiting consumer loop")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					logc.Info(ctx, "ReadMessage canceled")
					return
				}
				logc.Errorf(ctx, "ReadMessage error: %v", err)
				time.Sleep(time.Second)
				continue
			}
			msgCh <- msg
		}
	}
}

//func (c *Consumer) ConsumeMessagesWithContext(handler MessageHandle) {
//	go func() {
//		defer func() {
//			if err := recover(); err != nil {
//				logc.Errorf(context.Background(), "ConsumeMessagesWithContext recover error:%v, %s", err, string(debug.Stack()))
//			}
//		}()
//		defer func() {
//			logc.Error(context.Background(), "******************************************************  MQ consume quit ****************************************************** ")
//		}()
//		for {
//			newCtx := context.Background()
//			msg, err := c.reader.ReadMessage(newCtx)
//			if err != nil {
//				time.Sleep(time.Second)
//				continue
//			}
//			if msg.Value == nil || string(msg.Value) == "" {
//				continue
//			}
//
//			logc.Infof(newCtx, "---- kafka:ConsumeMessagesWithContext:topic: %s, msg: %s", msg.Topic, string(msg.Value))
//			ka := &KafkaMessage{}
//			err = json.Unmarshal(msg.Value, ka)
//			if err != nil {
//				continue
//			}
//
//			newCtx = ka.SetContext(newCtx)
//			for i := int64(1); i < 4; i++ { // 最大重试次数
//				err = handler(newCtx, ka.GetMsg())
//				if err == nil {
//					break
//				}
//				//if i == 3 {
//				//c.sendDeadLetterQueue(newCtx, msg.Topic, ka)
//				//break
//				//}
//				time.Sleep(time.Second) // 等待一段时间
//			}
//		}
//	}()
//	select {}
//}
