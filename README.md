# ls-having

A command-line tool for finding directories based on specified flag files and other conditions

[![ls-having](https://snapcraft.io/ls-having/badge.svg)](https://snapcraft.io/ls-having)
[![codecov](https://codecov.io/gh/handy-common-utils/ls-having/branch/master/graph/badge.svg?token=CJLY2DXUAU)](https://codecov.io/gh/handy-common-utils/ls-having)
[![Go Reference](https://pkg.go.dev/badge/github.com/handy-common-utils/ls-having.svg)](https://pkg.go.dev/github.com/handy-common-utils/ls-having/lsh)

*ls-having* is a powerful tool that makes it easy to find directories. With a simple and intuitive command-line interface, it allows you to quickly search for directories based on specified flag files and other conditions, saving you time and effort. Whether you're a developer who works with mono repos, or you just need an easy way to find your directories, *ls-having* is the perfect tool for you.

- Flexible search options: it allows you to specify various options and conditions for searching for directories, such as flag files, the check file, regular expressions for matching the check file, excludes, and maximum depth.
- Configurable error handling: there is an option for configuring how it handles errors such as non-existing directory or no access permission, allowing users to customize and integrate the tool for different situations.
- Cross-platform support: *ls-having* has native executables available for Linux, MacOS, and Windows.
- Small and quick: *ls-having* has just one native executable file needed, there is no dependency on anything else, and the tool is just about 2MB in size. It starts and runs very quickly.

## Installation

**MacOS (Intel or Apple silicon) and Linux and WSL**

On MacOS, you can use the [Homebrew](https://brew.sh/) package manager to install *ls-having*:

```shell
brew install handy-common-utils/tap/ls-having
```

Homebrew also supports Linux and WSL.

**Linux (all kinds of, except Ubuntu in WSL)**

On Linux, *ls-having* is available in the [snap store](https://snapcraft.io/ls-having) and can be installed through snap:

```shell
sudo snap install ls-having
```

Snap store is preinstalled in some Linux distributions (such like Ubuntu).
If it is not available, you can [install it](https://snapcraft.io/docs/installing-snapd) by yourself.

As of Janurary 2023, a component (`snapd`) required by snap does not work in Ubuntu running in WSL.
The workaround is to install through Homebrew.

**Manual download (works for Windows and others)**

For Windows, you can manually download the *ls-having* executable from the [releases page on GitHub](https://github.com/handy-common-utils/ls-having/releases) and unzip and copy the file to the desired location.

You can download executables and install for other operating systems as well.

## Quick start

To use *ls-having*, you can specify the root directory for the search as the last argument on the command line, or you can omit it to use the current directory as the root directory. The `-f`/`--flag-file` option is required to specify the name or glob of the flag file that the tool should look for. The `-c`/`--check-file` option and `-e`/`--check-regexp` option can be specified to refine the search.

For example, to search for directories containing a "package.json" file, you would run the following command:

```
ls-having -f package.json
```

This will search for directories containing a "package.json" file in the current directory and print the names of those directories to the console. You can then use the output of *ls-having* as arguments for other commands, such as `xargs`, to perform actions on those directories.

For more information on how to use *ls-having* and the available options and arguments, you can refer to the tool's [documentation on GitHub](https://github.com/handy-common-utils/ls-having) or run the `ls-having --help` command.

## Usage - as a CLI tool

### Help

`ls-having -h` prints the help screen:

```
Usage: ls-having -f name-or-glob [options] [root-dir]
Options:
  -c, --check-file name           name of the additional file to check
  -i, --check-inverse             regard regular expression not matching as positive
  -e, --check-regexp expression   regular expression for testing the content of the check file (default ".*")
  -d, --depth int                 how deep to look into subdirectories, 0 means only look at root directory, -1 means no limit (default 5)
  -r, --error ignore|panic|print  how (ignore|panic|print) to handle errors such like non-existing directory, no access permission, etc. (default "ignore")
  -x, --exclude glob              glob of the directories to exclude, this option can appear multiple times
  -f, --flag-file glob            name or glob of the flag file, this option can appear multiple times
  -h, --help                      show help information
  -n, --no-default-excludes       don't apply default excludes
  -0, --print0                    separate paths in the output with null characters (instead of newline characters)
  -s, --subdirectories-only       don't return root directory even if it meets conditions
References:
  Glob syntax: https://github.com/gobwas/glob#example
  Regexp syntax: https://pkg.go.dev/regexp/syntax
  Home page: https://github.com/handy-common-utils/ls-having
```

### Root directory

The root directory for the search, if not omitted, must be specified as the last argument in the command.

If the root directory is not specified, the current directory (./) will be used as the root directory for the search.

### Flag file

You must specify at least one flag file (`-f`/`--flag-file`) or check file (`-c`/`--check-file`),
otherwise *ls-having* would print out an error message and exit.

You can specify multiple flag files by using the `-f`/`--flag-file` option multiple times.
In such case, directories having **any** of those files or matching **any** of those globs will be returned.

If you have check file (`-c`/`--check-file`) specified and want to use the check file as the flag file, you can omit the `-f`/`--flag-file` option.
In this case, the check file will also be used as the flag file.

For example, to search for directories containing "package.json" or "build.gradle" or "build.gradle.kts" or "mvn.xml", you could run the following command:

```shell
ls-having -f package.json -f 'build.gradle*' -f mvn.xml
```

Please note that if you use `*` in the argument, you may need to quote the argument with single quotes,
otherwise the shell could interpret and translate it before it reaches the program.

### Check file

The check file can be useful for specifying additional conditions for the directories that *ls-having* searches for, allowing you to find more specific sets of directories.

The check file is specified using the `-c` or `--check-file` option, and its value can be a relative path to the directory.

The check file specified by `-c` or `--check-file` is an additional file
that can be checked in the directories having the flag file.
Its value can be a relative path.

For example, to search for directories containing a "serverless.*" file and also a "build.gradle" file, you could run the following command:

```shell
ls-having -f 'serverless.*' -c build.gradle
```

In addition to specifying the check file, you can also use the `-e` or `--check-regexp` option to specify a regular expression that the content of the check file must match in order for the directory to be included in the search results. If omitted, `.*` is used as the default check expression.

For example, to search for directories containing a "mvn.xml" file and also a "elastic-beanstalk.config" file containing string "MY_ENV_NAME:" in a subdirectory named ".ebextentions", you could run the following command:

```shell
ls-having -f mvn.xml -c .ebextentions/elastic-beanstalk.config -e 'MY_ENV_NAME:'
```

### Default excludes

By default, these directories are not looked into:

- `.git` and `**/.git`
- `node_modules` and `**/node_modules`
- `testdata` and `**/testdata`

Option `-n`/`--no-default-excludes` can be used to disable this behaviour.

Option `-x`/`--exclude` can be used to add more globs to the list.
This Option can appear multiple times.

Examples, to replace "testdata" with "fixtures" in the list of excludes,
you would do this:

- `ls-having -f package.json --no-default-excludes --exclude node_modules --exclude '**/node_modules' --exclude .git --exclude '**/.git' --exclude fixtures --exclude '**/fixtures'`.

### Default maximum depth

By default, *ls-having* will only search for directories up to 5 levels deep in the directory tree. The root directory is considered as level 0, and its direct subdirectories are considered as level 1, and so on.

If you want to search for directories at a deeper level, you can use the `--depth` or `-d` option to specify the maximum depth that *ls-having* should search. For example, if you specify `--depth 10`, the tool will search for directories up to 10 levels deep in the directory tree.

If you specify a negative number for the `--depth` option, *ls-having* will continue searching for directories until there are no more subdirectories to search. This can be useful for searching the entire directory tree without having to specify a specific maximum depth.

### Error handling

By default, *ls-having* will not print out any error messages if it encounters errors during processing. For example, if the root directory does not exist, or if the user does not have permission to access the root directory or any subdirectories, *ls-having* will simply continue processing without printing out any error messages, and then exit with code `0`.

However, you can use the `--error` option to change how *ls-having* handles errors. If you specify the `--error panic` option, if it encounters any errors during processing, *ls-having* will stop immediately, print out an error message to `stderr`, and exit with code `1`. This can be useful if you want to fail CI/CD pipeline immediately.

If you specify the `--error print` option, *ls-having* will ignore errors during processing and exit with code `0`, but it will also print out all error messages to `stderr`. This can be useful if you want to continue processing because you know those errors are caused by directory access permission issues.

### Examples

‣ Find all directories in `./` having `package.json` file,
and run `npm audit fix` in them one by one:

```shell
ls-having -f package.json | xargs -I {} bash -c 'cd "{}"; npm audit fix'
```

‣ Find all directories in `./` having `package.json` file
and the `package.json` file does not contain `"volta":`:

```shell
ls-having -c package.json -i -e '"volta":'
```

‣ Find all directories in `./` having `package.json` file
and the `package.json` file has `mocha` specified as a dependency,
then for each of those directories reinstall latest version of `mocha` as dev-dependency:

```shell
ls-having -c package.json -e '"dependencies":\s*{[^{}]*"mocha":' | xargs -I {} bash -c 'cd {}; npm i -D mocha@latest'
```

‣ Find all directories in `./` having `package.json` file,
go as deep as 8 levels, and don't apply default excludes
(such like `.git` and `node_modules`):

```shell
ls-having -f package.json -d 8 -n
```

‣ Find all directories in `testdata/repo1` having `serverless.*`
(such like `serverless.yml`, `serverless.ts`, `serverless.js`),
and also having `package.json`:

```shell
ls-having  -f 'serverless.*' -c package.json testdata/repo1
```

Please note that if you use `*` in the argument as part of a
glob or regular expression, you need to quote the argument with single quotes,
otherwise the shell could interpret and translate it before it reaches the program.

‣ Find all directories having `cdk.json`,
and also a `package.json` file containing text `"@types/mocha":`:

```shell
ls-having -f cdk.json -c package.json -e '"@types/mocha":'
```

‣ Find all directories having `package.json` and also a `node_modules/package1/package.json` file in its subdirectory:

```shell
ls-having -f package.json -c node_modules/package1/package.json
```

Please note that the check file used above is a relative path,
and the excluding logic does not apply to this path.
You can even use `..` and the check "file" could be a directory,
below are more examples:

```shell
ls-having -f package.json -c ../australia/serverless.yml -e datadog
ls-having -f package.json -c ../australia
ls-having -f package.json -c ../australia -i
```

‣ Find all subdirectories (but exclude the current directory `./`)
having `package.json`,
and also having `serverless.yml` file contain text `datadog`:

```shell
ls-having -f package.json -c serverless.yml -e 'datadog' -s
```

‣ Find all subdirectories under `/tmp/sample/repo` (but exclude the root directory `/tmp/sample/repo`)
having `build.gradle*` or `mvn.xml`:

```shell
ls-having -f 'build.gradle*' -f 'mvn.xml' -s /tmp/sample/repo
```

‣ Find all directories in `/tmp/sample/repo`
having `build.gradle*` and print out their details:

```shell
ls-having -f 'build.gradle*' /tmp/sample/repo | xargs -I {} bash -c 'cd {}; pwd; ls -l build.gradle*'
```

‣ Make sure all the node versions specified by `.nvmrc` files are installed:

```shell
ls-having -f .nvmrc | xargs -I {} bash -c '. ~/.nvm/nvm.sh; cd "{}"; nvm install'
```

Please note that sourcing `~/.nvm/nvm.sh` is needed because nvm alias is not automatically avaible in subshells.

‣ List all different node versions specified in `.nvmrc` files:

```shell
ls-having -f .nvmrc | xargs -I {} bash -c 'cat {}/.nvmrc; echo;' | sed 's/^v// ; /^$/d' | sort -u
```

‣ Install dependencies in those directories having `.nvmrc` but not using `yarn`:

```shell
ls-having -f .nvmrc -c yarn.lock -i | xargs -I {} bash -c '. ~/.nvm/nvm.sh; cd "{}"; nvm use; npm ci'
```

‣ Install dependencies in those directories having `.nvmrc` and using `yarn`:

```shell
ls-having -f .nvmrc -c yarn.lock | xargs -I {} bash -c '. ~/.nvm/nvm.sh; cd "{}"; nvm use; npm i -g yarn; yarn install --frozen-lockfile'
```

‣ Use null character (instead of newline character) as separator in the output,
so that `-0` option of xargs can be used:

```shell
ls-having --print0 -f package.json testdata/repo1 | xargs -0 -I {} bash -c 'cd "{}"; pwd'
```

## Usage - as a Go package

Package summary page: https://pkg.go.dev/github.com/handy-common-utils/ls-having/lsh

```go
import (
	"fmt"
	"github.com/handy-common-utils/ls-having/lsh"
)

func main() {
	var dirs, errors = lsh.LsHaving(&lsh.Options{
		// options here
	}, rootDir)
	fmt.Println(dirs)
}
```

See [main.go](https://github.com/handy-common-utils/ls-having/blob/master/main.go#:~:text=lsh.LsHaving) for code example.

## Contributing

**Run locally**

```
go run . <command line arguments>
```

**Test (with coverage percentage shown)**

```
go test -coverpkg=github.com/handy-common-utils/ls-having/lsh,github.com/handy-common-utils/ls-having
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
