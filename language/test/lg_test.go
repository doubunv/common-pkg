package test

import (
	"fmt"
	"github.com/doubunv/common-pkg/language"
	"testing"
)

func TestIsExist(t *testing.T) {
	// 示例数据
	data := map[string]interface{}{
		"foo": "Hello",
		"bar": []interface{}{
			"apple",
			map[string]interface{}{
				"nested": "banana",
			},
			"cherry",
		},
	}

	fmt.Println("Before modification:", data)
	res := language.SwitchLanguage(data, "india")
	fmt.Println("After modification:", res)
}
