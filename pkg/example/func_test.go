package example_test

import (
	"testing"

	"github.com/9506hqwy/template-go-module/pkg/example"
)

func Test_add(t *testing.T) {
	ret := example.Add(1, 2)
	if ret != 3 {
		t.Error("Failed.", ret)
	}
}

func Test_sub(t *testing.T) {
	ret := example.Sub(2, 1)
	if ret != 1 {
		t.Error("Failed.", ret)
	}
}
