package main

import (
	"fmt"
)

func main() {
	test := make(map[string]int)
	test["a"] = 1
	test["b"] = 2
	fmt.Println(len(test))
}
