package v1

import (
	"context"
	"fmt"
	"github.com/doubunv/common-pkg/amqp/kafka/config"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logc"
	"time"
)

type Producer struct {
	Writer *kafka.Writer
	config config.ProviderConfig
}

func NewProducer(config config.ProviderConfig) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  config.Topic,
		Balancer:               &kafka.LeastBytes{},
		MaxAttempts:            5,               // 最大重试次数
		WriteTimeout:           3 * time.Second, // 写入超时时间
		AllowAutoTopicCreation: true,
		RequiredAcks:           config.RequiredAcks,
	}
	return &Producer{
		Writer: writer,
		config: config,
	}
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}

func (p *Producer) ProduceMessageWithContext(ctx context.Context, message *KafkaMessage) error {
	logc.Info(ctx, fmt.Sprintf("--- kafka:ProduceMessage: topic:%s, Message:%s", p.config.Topic, string(message.PacketMsg(p.config.Topic).Value)))
	err := p.Writer.WriteMessages(ctx, message.PacketMsg(p.config.Topic))
	if err != nil {
		logc.Error(ctx, fmt.Sprintf("--- kafka:ProduceMessageWithContext: topic:%s, Message:%s", p.config.Topic, string(message.PacketMsg(p.config.Topic).Value)))
		return err
	}

	return nil
}
