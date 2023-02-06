package utils

import (
	"fmt"
	"testing"
)

func TestPointer(t *testing.T) {
	a := []*string{}
	for i := 0; i < 10; i++ {
		value := "a"
		a = append(a, &value)
	}
}

func TestReverseSlice(t *testing.T) {
	a := []float64{11, 33, 44, 55, 66, 33, 22}

	ReverseSlice(a)
	fmt.Println(a)
}

func TestGetUUID(t *testing.T) {
	str := GetUUID()
	fmt.Println(str)
}
