package v1

import (
	"context"
	"errors"
	"github.com/doubunv/common-pkg/amqp/kafka/config"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

type MessageHandler func(ctx context.Context, msg string) error

type Consumer struct {
	reader    *kafka.Reader
	workers   map[int]*partitionWorker
	mu        sync.Mutex
	maxRetry  int
	kafkaConf kafka.ReaderConfig
}

func NewConsumer(conf config.CustomerConfig) *Consumer {
	kafkaConf := kafka.ReaderConfig{
		Brokers:     conf.Brokers,
		GroupID:     conf.GroupID,
		Topic:       conf.Topic,
		StartOffset: conf.StartOffset,
	}
	reader := kafka.NewReader(kafkaConf)
	return &Consumer{
		reader:    reader,
		workers:   make(map[int]*partitionWorker),
		maxRetry:  15,
		kafkaConf: kafkaConf,
	}
}

func (c *Consumer) ConsumeMessagesWithContext(ctx context.Context, handler MessageHandler) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("consumer stopped")
				return
			}
			time.Sleep(time.Second)
			continue
		}

		worker := c.getWorker(ctx, msg.Partition, handler)
		worker.ch <- msg
	}
}

func (c *Consumer) getWorker(ctx context.Context, partition int, handler MessageHandler) *partitionWorker {

	c.mu.Lock()
	defer c.mu.Unlock()

	if w, ok := c.workers[partition]; ok {
		return w
	}

	w := newPartitionWorker(partition, c.reader, c.maxRetry, handler)
	c.workers[partition] = w
	go w.run(ctx)

	return w
}
