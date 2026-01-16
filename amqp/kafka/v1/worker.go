package v1

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logc"
	"log"
	"runtime/debug"
	"time"
)

type partitionWorker struct {
	partition int
	ch        chan kafka.Message
	reader    *kafka.Reader
	maxRetry  int
	handler   MessageHandler
}

func newPartitionWorker(partition int, reader *kafka.Reader, maxRetry int, handler MessageHandler) *partitionWorker {

	return &partitionWorker{
		partition: partition,
		ch:        make(chan kafka.Message, 100),
		reader:    reader,
		maxRetry:  maxRetry,
		handler:   handler,
	}
}

func (w *partitionWorker) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("partition %d worker stopped\n", w.partition)
			return

		case msg := <-w.ch:
			w.process(ctx, msg)
		}
	}
}

func (w *partitionWorker) process(ctx context.Context, msg kafka.Message) {
	defer func() {
		if r := recover(); r != nil {
			logc.Errorf(ctx, "panic: %v\n%s\n", r, debug.Stack())
			sendDLQ(ctx, msg)
			_ = w.reader.CommitMessages(ctx, msg)
		}
	}()

	var km KafkaMessage
	if err := json.Unmarshal(msg.Value, &km); err != nil {
		log.Printf("unmarshal error: %v\n", err)
		_ = w.reader.CommitMessages(ctx, msg)
		return
	}

	msgCtx := km.SetContext(ctx)

	var err error
	for i := 1; i <= w.maxRetry; i++ {
		if err = w.handler(msgCtx, km.GetMsg()); err == nil {
			_ = w.reader.CommitMessages(ctx, msg)
			return
		}
		time.Sleep(time.Duration(i) * time.Second)
	}

	sendDLQ(ctx, msg)
	_ = w.reader.CommitMessages(ctx, msg)
}
