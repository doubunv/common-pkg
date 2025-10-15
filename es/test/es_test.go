package test

import (
	"context"
	"fmt"
	"github.com/doubunv/common-pkg/es/esaws/core"
	"github.com/doubunv/common-pkg/es/esaws/model"
	"strconv"
	"testing"
	"time"
)

// 定义一个实现 IndexTable 接口的结构体
type MyIndexTable struct {
	_indexName string
	_id        string
	UserId     int64   `json:"user_id"`
	PUserId    []int64 `json:"p_user_id"`
	Amount     float64 `json:"amount"` //es在做map映射的时候需要用double，不然会有精度丢失
}

// 实现 IndexName 方法
func (m MyIndexTable) IndexName() string {
	return m._indexName
}

// 实现 GetId 方法
func (m MyIndexTable) GetId() string {
	return m._id
}

// 实现 SetId 方法
func (m *MyIndexTable) SetId(id string) {
	m._id = id
}

func TestLogInfo(t *testing.T) {

	esClient := core.MustNewEs(&core.Config{
		Addresses:  []string{"https://search-es-usdt-db-dev-fp3bbysw7vfmnn6u7segtwqxpa.aos.ap-southeast-1.on.aws"},
		Username:   "admin",
		Password:   "Dev123456.",
		MaxRetries: 3,
	})
	ctx := context.Background()
	esModel := model.NewEsModel(ctx, esClient)

	//创建
	tm := time.Now().Unix()
	data := &MyIndexTable{_indexName: "example_index", UserId: 1759848897, Amount: 1.78, PUserId: []int64{tm, tm, tm + 2, tm + 3}, _id: strconv.FormatInt(tm, 10)}
	err := esModel.InsertSchema(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("end")

	//查询单条
	//data1 := &MyIndexTable{Name: "example_index"}
	//err := esModel.FindOne("1111", data1)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(data1)

	//修改
	//data2 := &MyIndexTable{_indexName: "example_index", UserId: time.Now().Unix(), _id: "1111"}
	//err := esModel.UpdateSchema(data2)
	//if err != nil {
	//	return
	//}

	//基本查询
	//data3 := &MyIndexTable{_indexName: "example_index", _id: "1111"}
	//resData := make([]MyIndexTable, 0)
	//res, num, err := esModel.Search(data3, &resData, map[string]interface{}{
	//	"from": 0, //todo 是从0开始的哟
	//	"size": 2,
	//	"query": map[string]interface{}{
	//		"match": map[string]interface{}{ //todo 精确查询
	//			"user_id": 1759842237,
	//		},
	//		"term": map[string]interface{}{ //todo 查询数组
	//			"p_user_id": 1759842237,
	//		},
	//		"terms": map[string]interface{}{ //todo 查询多个
	//			"p_user_id": []string{"1759842237", "1759842238"},
	//		},
	//	},
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(res, num)

	//简单聚合查询
	//data4 := &MyIndexTable{_indexName: "example_index", _id: "1111"}
	//resData := MyIndexTable{}
	//res, num, aggregate, err := esModel.Search(data4, &resData, map[string]interface{}{
	//	"size": 0, // 不需要返回具体文档
	//	"aggs": map[string]interface{}{
	//		"user_id": map[string]interface{}{
	//			"sum": map[string]interface{}{
	//				"field": "user_id", // 要累加的字段
	//			},
	//		},
	//	},
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(aggregate)
	//fmt.Println(res, num)

	//分组聚合查询
	//data4 := &MyIndexTable{_indexName: "example_index", _id: "1111"}
	//resData := MyIndexTable{}
	//res, num, aggregate, err := esModel.Search(data4, &resData, map[string]interface{}{
	//	"size": 0, // 不需要返回具体文档
	//	"aggs": map[string]interface{}{
	//		"amount_by_user": map[string]interface{}{
	//			"terms": map[string]interface{}{
	//				"field": "user_id", // 按 user_id 分组
	//			},
	//			"aggs": map[string]interface{}{ // 对每组进行求和
	//				"total_amount": map[string]interface{}{
	//					"sum": map[string]interface{}{
	//						"field": "amount", // 要累加的字段
	//					},
	//				},
	//			},
	//		},
	//	},
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(aggregate)
	//fmt.Println(res, num)
}
