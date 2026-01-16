package v2

import (
	"context"
	"errors"
	"github.com/doubunv/common-pkg/amqp/kafka/config"
	"testing"
	"time"
)

func Test_Comsume(t *testing.T) {
	cf := config.CustomerConfig{
		ProviderConfig: config.ProviderConfig{
			Brokers: []string{"xx.xx.xx.xx:9092"},
			Topic:   "xxxx",
		},
		GroupID: "lala",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cf = config.CustomerConfig{
		ProviderConfig: config.ProviderConfig{
			Brokers: []string{"xx.xx.xx.xx:9092"},
			Topic:   "xxxx",
		},
		GroupID: "lala",
	}
	consumer := NewConsumer(cf)

	//返回 nil 代表消费成功
	go consumer.Consume(ctx, func(ctx context.Context, msg string) error {
		//log.Printf("consume msg: %s\n", msg)

		// 模拟失败
		if msg == "fail" {
			return errors.New("mock error")
		}
		time.Sleep(200 * time.Millisecond)
		return nil
	})
}
