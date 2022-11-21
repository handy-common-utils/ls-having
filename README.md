# ls-having
List (or operate in) directories having specified flag files and/or meeting other specified conditions 

## Usage - as a CLI tool

## Usage - as a Go package

```go
import (
	"fmt"
	"github.com/handy-common-utils/ls-having/lsh"
)

func main() {
	var dirs = lsh.LsHaving()
	fmt.Println(dirs)
}

```
