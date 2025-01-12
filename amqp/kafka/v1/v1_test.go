package v1

import (
	"context"
	"github.com/doubunv/common-pkg/amqp/kafka/config"
	"github.com/doubunv/common-pkg/headInfo"
	"testing"
)

func Test_Comsume(t *testing.T) {
	cf := config.CustomerConfig{
		ProviderConfig: config.ProviderConfig{
			Brokers: []string{"54.179.172.191:9092"},
			Topic:   "wallet_balance_log_topic",
		},
		GroupID: "lala",
	}

	//返回 nil 代表消费成功
	NewConsumer(cf).ConsumeMessagesWithContext(func(ctx context.Context, msg string) error {
		//fmt.Println(headInfo.GetBusinessCode(ctx))
		//fmt.Println(strconv.FormatInt(headInfo.GetTokenUid(ctx), 10))
		//fmt.Println(headInfo.GetTrance(ctx))
		//fmt.Println(msg)

		return nil
	})
}

func Test_Provider(t *testing.T) {
	cf := config.ProviderConfig{
		Brokers: []string{"54.179.172.191:9092"},
		Topic:   "wallet_balance_log_topic",
	}

	ctx := headInfo.SetBusinessCode(context.Background(), "21111111")

	km := NewKafkaMessage(ctx, "ffdjfdfefefeefewtest(*******")
	NewProducer(cf).ProduceMessageWithContext(ctx, km)
}
