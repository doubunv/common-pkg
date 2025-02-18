package v1

import (
	"context"
	"fmt"
	"github.com/doubunv/common-pkg/amqp/kafka/config"
	"github.com/doubunv/common-pkg/headInfo"
	"math/rand"
	"testing"
)

func Test_Comsume(t *testing.T) {
	cf := config.CustomerConfig{
		ProviderConfig: config.ProviderConfig{
			Brokers: []string{"xx.xx.xx.xx:9092"},
			Topic:   "xxxx",
		},
		GroupID: "lala",
	}

	//返回 nil 代表消费成功
	NewConsumer(cf).ConsumeMessagesWithContext(func(ctx context.Context, msg string) error {
		//fmt.Println(headInfo.GetBusinessCode(ctx))
		//fmt.Println(strconv.FormatInt(headInfo.GetTokenUid(ctx), 10))
		//fmt.Println(headInfo.GetTrance(ctx))
		//fmt.Println(msg)

		randomNum := rand.Intn(5) + 1
		if randomNum == 1 {
			return nil
		}

		return fmt.Errorf("处理失败")
		//return errors.New("aaaa")
	})
}

func Test_Provider(t *testing.T) {
	cf := config.ProviderConfig{
		Brokers: []string{"xx.xx.xx.xx:9092"},
		Topic:   "xxxx",
	}

	ctx := headInfo.SetBusinessCode(context.Background(), "21111111")

	km := NewKafkaMessage(ctx, "ffdjfdfefefeefewtest(*******")
	NewProducer(cf).ProduceMessageWithContext(ctx, km)
}
