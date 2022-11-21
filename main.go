package main

import (
	"fmt"

	"github.com/handy-common-utils/ls-having/lsh"
)

func main() {
	fmt.Println("Hello!")
	var dirs = lsh.LsHaving()
	fmt.Println(dirs)
}
