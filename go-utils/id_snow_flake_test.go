package goutils

import (
	"fmt"
	"testing"
)

func TestGenId(t *testing.T) {
	//GenIdInit(&SnowFlakeId{machineId: 2})
	//
	//for i := 0; i < 5000; i++ {
	//	fmt.Println(GenId())
	//}

	aORb := IfNot[[]string](true, []string{"aa"}, []string{})
	fmt.Print(aORb)
}
