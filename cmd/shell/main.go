package main

import (
	"fmt"

	"github.com/9506hqwy/template-go-module/pkg/example"
)

func main() {
	ret := example.Add(2, 5)
	fmt.Printf("result = %d\n", ret)
}
