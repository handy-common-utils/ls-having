# ls-having

This is a tool that can list directories having specified flag files and optionally meeting other specified conditions.

[![ls-having](https://snapcraft.io/ls-having/badge.svg)](https://snapcraft.io/ls-having)

## Install

**MacOS (Intel or Apple silicon)**

`ls-having` can be installed through [Homebrew](https://brew.sh/):

```sh
brew install handy-common-utils/tap/ls-having
```

To upgrade:

```sh
brew upgrade handy-common-utils/tap/ls-having
```

**Linux (all kinds of)**

`ls-having` can be installed through [snap](https://snapcraft.io/docs/installing-snapd):

```sh
sudo snap install ls-having
```

To upgrade:

```sh
sudo snap refresh ls-having
```

## Usage - as a CLI tool

### Help

```
Usage: ls-having -f name-or-glob [options] [root-dir]
Options:
  -c, --check-file name          name of the additional file to check
  -i, --check-inverse            regard regular expression not matching as positive
  -e, --check-regexp expression  regular expression for testing the content of the check file (default ".*")
  -d, --depth int                how deep to look into subdirectories, 0 means only look at root directory, -1 means no limit (default 5)
  -x, --exclude glob             glob of the directories to exclude, this option can appear multiple times
  -f, --flag-file glob           name or glob of the flag file, this option can appear multiple times
  -h, --help                     show help information
  -n, --no-default-excludes      don't apply default excludes
  -s, --subdirectories-only      don't return root directory even if it meets conditions
References:
  Glob syntax: https://github.com/gobwas/glob#example
  Regexp syntax: https://pkg.go.dev/regexp/syntax
  Home page: https://github.com/handy-common-utils/ls-having
```

### Examples

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