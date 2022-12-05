# ls-having

This is a tool that can list directories having specified flag files and optionally meeting other specified conditions.

## Install

**MacOS (Intel or Apple silicon)**

`ls-having` can be installed through [Homebrew](https://brew.sh/):

```sh
brew install handy-common-utils/tap/ls-having
```

**Linux (all kinds of)**

`ls-having` can be installed through [snap](https://snapcraft.io/docs/installing-snapd):

```sh
sudo snap install ls-having
```

## Usage - as a CLI tool

Find all directories in `./` having `package.json` file,
go as deep as 8 levels, and don't apply default excludes
(such like `.git` and `node_modules`):

```sh
ls-having -f package.json -d 8 -n
```

Find all directories in `testdata/repo1` having `serverless.*`
(such like `serverless.yml`, `serverless.ts`, `serverless.js`),
and also having `package.json`:

```sh
ls-having  -f 'serverless.*' -c package.json testdata/repo1
```

Find all directories in `./` having `package.json`,
and the `package.json` file must contain text `mocha`:

```sh
ls-having -f 'package.json' -c package.json -e 'mocha'
```

Find all subdirectories under `./` (the root directory `./` is excluded)
having `package.json`,
and also having `serverless.yml` file contain text `datadog`:

```sh
ls-having -f 'package.json' -c serverless.yml -e 'datadog' -s
```

Find all subdirectories under `/tmp/sample/repo` (the root directory `./` is excluded)
having `build.gradle*` or `mvn.xml`:

```sh
ls-having -f 'build.gradle*' -f 'mvn.xml' -s /tmp/sample/repo
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