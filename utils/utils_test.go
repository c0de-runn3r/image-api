package utils

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

// func TestOriginalSize(t *testing.T) {
// 	width, height, _ := originalSize("cat.png")
// 	fmt.Printf("width: %d, height: %d\n", width, height)
// }
