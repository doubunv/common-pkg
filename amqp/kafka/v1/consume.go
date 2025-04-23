package v1

import (
	"context"
	"encoding/json"
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
		Brokers: conf.Brokers,
		GroupID: conf.GroupID,
		Topic:   conf.Topic,
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

func (c *Consumer) ConsumeMessagesWithContext(handler MessageHandle) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logc.Errorf(context.Background(), "ConsumeMessagesWithContext error:%v, %s", err, string(debug.Stack()))
			}
		}()
		for {
			newCtx := context.Background()
			msg, err := c.reader.ReadMessage(newCtx)
			logc.Infof(newCtx, "---- kafka:ConsumeMessagesWithContext:topic: %s, msg: %s", msg.Topic, string(msg.Value))
			if string(msg.Value) == "" {
				return
			}

			ka := &KafkaMessage{}
			err = json.Unmarshal(msg.Value, ka)
			if err != nil {
				logc.Errorf(newCtx, "---- kafka:ConsumeMessagesWithContext:topic: %s, msg:%+v, err: %+v", msg.Topic, string(msg.Value), err)
				continue
			}
			newCtx = ka.SetContext(newCtx)
			go func(msg kafka.Message) {
				defer func() {
					if err := recover(); err != nil {
						logc.Errorf(context.Background(), "ConsumeMessagesWithContext handler error:%v, %s, %s", string(msg.Value), err, string(debug.Stack()))
					}
				}()
				for i := 1; i < 4; i++ { // 最大重试次数
					err = handler(newCtx, ka.GetMsg())
					if err == nil {
						break
					}
					if i == 2 {
						c.sendDeadLetterQueue(newCtx, msg.Topic, ka)
						break
					}
					time.Sleep(500 * time.Millisecond) // 等待一段时间
				}
			}(msg)
		}
	}()
	select {}
}
