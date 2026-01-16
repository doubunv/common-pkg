package v1

import (
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logc"
)

import (
	"context"
)

func sendDLQ(ctx context.Context, msg kafka.Message) {
	// 这里可以换成真正的 DLQ producer
	logc.Errorf(
		ctx,
		"[DLQ] topic=%s partition=%d offset=%d value=%s\n",
		msg.Topic,
		msg.Partition,
		msg.Offset,
		string(msg.Value),
	)
}
