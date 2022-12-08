# ls-having

This is a tool that can list directories having specified flag files and optionally meeting other specified conditions.

[![ls-having](https://snapcraft.io/ls-having/badge.svg)](https://snapcraft.io/ls-having)
[![codecov](https://codecov.io/gh/handy-common-utils/ls-having/branch/master/graph/badge.svg?token=CJLY2DXUAU)](https://codecov.io/gh/handy-common-utils/ls-having)
[![Go Reference](https://pkg.go.dev/badge/github.com/handy-common-utils/ls-having.svg)](https://pkg.go.dev/github.com/handy-common-utils/ls-having/lsh)

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

**Manual download (works for Windows)**

You can just download, unzip, and copy the executable to anywhere you like:
https://github.com/handy-common-utils/ls-having/releases

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

### Default excludes

By default, these directories will not be looked into:

- `.git` and `**/.git`
- `node_modules` and `**/node_modules`
- `testdata` and `**/testdata`

Flag `--no-default-excludes` can be used to disable this behaviour.

Flag `--exclude` can be used to add more globs to the list.
This flag can appear multiple times.

Examples:

- `ls-having -f package.json --no-default-excludes --exclude node_modules --exclude '**/node_modules' --exclude '**/sample'`.

### Examples

Find all subdirectories having `package.json` file,
and run `npm audit fix` in them one by one:

```sh
ls-having -f package.json | xargs -I {} npm audit fix
```

The `-I {}` flag above has very similar effect as `-L 1` flag.


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

Please note that if you use `*` in the argument as part of a
glob or regular expression, you need to quote the argument with single quotes,
otherwise the shell could interpret and translate it before it reaches the program.

Find all subdirectories having `package.json`,
and the `package.json` file must contain text `"@types/mocha":`:

```sh
ls-having -f 'package.json' -c package.json -e '"@types/mocha":'
```

Find all subdirectories (the current directory `./` is excluded)
having `package.json`,
and also having `serverless.yml` file contain text `datadog`:

```sh
ls-having -f 'package.json' -c serverless.yml -e 'datadog' -s
```

Find all subdirectories under `/tmp/sample/repo` (the root directory `/tmp/sample/repo` is excluded)
having `build.gradle` or `build.gradle.*` or `mvn.xml`:

```sh
ls-having -f 'build.gradle' -f 'build.gradle.*' -f 'mvn.xml' -s /tmp/sample/repo
```

## Usage - as a Go package

Package summary page: https://pkg.go.dev/github.com/handy-common-utils/ls-having/lsh

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

See [main.go](https://github.com/handy-common-utils/ls-having/blob/master/main.go#:~:text=var%20dirs%20%3D-,lsh.LsHaving,-(%26options) for example.

## Contributing

**Run locally**

```
go run . <command line arguments>
```

**Test**

```
go test
```

**Release**

1. Push to `master` branch
2. Tag with version number without prefix `v`.
   For example, use tag `1.2.34`.
3. Push the tag
4. GitHub workflow will automatically release to
   [snapcraft](https://snapcraft.io/ls-having),
   Homebrew,
   and [pkg.go.dev](https://pkg.go.dev/github.com/handy-common-utils/ls-having).
5. The workflow also automatically creates a tag with prefixed version number.
   For example, `v1.2.34`