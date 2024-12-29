package v1

import (
	"context"
	"encoding/json"
	"github.com/doubunv/common-pkg/headInfo"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/metadata"
)

type KafkaMessage struct {
	Msg   string
	HeadI map[string]interface{}
	Head  headInfo.Head
}

func NewKafkaMessage(ctx context.Context, msg string) *KafkaMessage {
	var res = &KafkaMessage{
		Msg:   msg,
		HeadI: make(map[string]interface{}),
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		for i, v := range md {
			if len(v) > 0 {
				res.HeadI[i] = v[0]
			}
		}
	}
	res.HeadI["trace"] = headInfo.GetTrance(ctx)
	return res
}

func (k *KafkaMessage) PacketMsg(key string) kafka.Message {
	pk := make(map[string]interface{})
	pk["msg"] = k.Msg
	pk["head"] = k.HeadI
	val, _ := json.Marshal(pk)
	return kafka.Message{
		Value: val,
		Key:   []byte(key),
	}
}

func (k *KafkaMessage) GetMsg() string {
	return k.Msg
}

func (k *KafkaMessage) SetContext(ctx context.Context) context.Context {
	newCtx := headInfo.ContextHeadInLog(ctx, &k.Head)
	return headInfo.HeadInMetadata(newCtx, k.Head)
}
