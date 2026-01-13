package commonTool

import (
	"fmt"
	"testing"
)

func Test_Comsume(t *testing.T) {

	//ll := GetZoneTime(0)
	//l8 := GetZoneTime(8)
	//fmt.Println(ll, l8)

	v7, err := NewUUidV7()
	if err != nil {
		return
	}
	fmt.Println(v7)
	fmt.Println(len(v7))
}
