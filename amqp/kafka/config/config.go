package config

import (
	"time"
)

type ConnectConfig struct {
	Network          string        `json:"network"`       // 网络协议
	Address          string        `json:"address"`       // 链接地址
	Topic            string        `json:"topic"`         // 主题
	Partition        int           `json:"partition"`     // 分区
	SendDeadlineTime time.Duration `json:"deadline_time"` // 发送超时时间
}

type CustomerConfig struct {
	ProviderConfig

	GroupID     string
	GroupTopics []string
	StartOffset int64 `json:"StartOffset,default=-1"` //LastOffset  int64 = -1  or FirstOffset int64 = -2
	Partition   int   `json:"partition,default=10"`
}

type ProviderConfig struct {
	Brokers []string
	Topic   string
}
