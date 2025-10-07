package test

import (
	"context"
	"fmt"
	"github.com/doubunv/common-pkg/es/esv7/core"
	"github.com/doubunv/common-pkg/es/esv7/model"
	"testing"
)

// 定义一个实现 IndexTable 接口的结构体
type MyIndexTable struct {
	name   string
	userId int64
	id     string
}

// 实现 IndexName 方法
func (m MyIndexTable) IndexName() string {
	return m.name
}

// 实现 GetId 方法
func (m MyIndexTable) GetId() string {
	return m.id
}

// 实现 SetId 方法
func (m *MyIndexTable) SetId(id string) {
	m.id = id
}

func TestLogInfo(t *testing.T) {
	esClient := core.MustNewEs(&core.Config{
		Addresses:  []string{"https://xxxx"},
		Username:   "admin",
		Password:   "",
		MaxRetries: 3,
	})

	ctx := context.Background()
	data := &MyIndexTable{name: "example_index", userId: 123, id: "1111"}
	err := model.NewEsModel(ctx, esClient).InsertSchema(data)
	if err != nil {
		fmt.Println(err)
		return
	}

}
