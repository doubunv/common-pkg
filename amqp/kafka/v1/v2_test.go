package v1

import (
	"context"
	"gitlab.888bbm.com/go-package/common-pkg/amqp/kafka/config"
	"gitlab.888bbm.com/go-package/common-pkg/headInfo"
	"testing"
)

func Test_Comsume(t *testing.T) {
	cf := config.CustomerConfig{
		ProviderConfig: config.ProviderConfig{
			Brokers: []string{"52.77.42.5:19092"},
			Topic:   "xq_topic",
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
		Brokers: []string{"52.77.42.5:19092"},
		Topic:   "xq_topic",
	}

	ctx := headInfo.SetBusinessCode(context.Background(), "21111111")

	km := NewKafkaMessage(ctx, "ffdjfdfefefeefewtest(*******")
	NewProducer(cf).ProduceMessageWithContext(ctx, km)
}
