package main

import (
	"fmt"

	"github.com/handy-common-utils/ls-having/lsh"
)

func main() {
	var dirs = lsh.LsHaving(nil, ".")
	fmt.Println(dirs)
}
