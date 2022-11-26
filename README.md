# ls-having
List (or operate in) directories having specified flag files and/or meeting other specified conditions 

## Usage - as a CLI tool

Find all directories in `./` having `package.json` file,
go as deep as 8 levels, and don't apply default excludes
(such like `.git` and `node_modules`):

```
ls-having -f package.json -d 8 -n
```

Find all directories in `testdata/repo1` having `serverless.*`
(such like `serverless.yml`, `serverless.ts`, `serverless.js`),
and also having `package.json`:

```
ls-having  -f 'serverless.*' -c package.json testdata/repo1
```

Find all directories in `./` having `package.json`,
and the `package.json` file must contain text `mocha`:

```
ls-having -f 'package.json' -c package.json -e 'mocha'
```

Find all subdirectories under `./` (the root directory `./` is excluded)
having `package.json`,
and also having `serverless.yml` file contain text `datadog`:

```
ls-having -f 'package.json' -c serverless.yml -e 'datadog' -s
```

## Usage - as a Go package

```go
import (
	"fmt"
	"github.com/handy-common-utils/ls-having/lsh"
)

func main() {
	var dirs = lsh.LsHaving(&lsh.Options{
		// options here
	}, rootDir)
	fmt.Println(dirs)
}
```

See [main.go](main.go) for example.