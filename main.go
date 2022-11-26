package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
	"github.com/handy-common-utils/ls-having/lsh"
	"github.com/josephvusich/go-getopt"
)

func main() {
	var optHelp = flag.Bool("help", false, "show help information")
	var optDepth = flag.Int("depth", 5, "how deep to look into subdirectories, 0 means only look at root directory, -1 means no limit")
	var optFlagFile = flag.String("flag-file", "", "name/glob of the flag file")
	var optCheckFile = flag.String("check-file", "", "name of the additional file to check")
	var optCheckRegexp = flag.String("check-regexp", ".*", "regular expression for testing the content of the check file")
	var optCheckInverse = flag.Bool("check-inverse", false, "regard regular expression not matching as positive")
	var optExcludes arrayFlag
	flag.Var(&optExcludes, "exclude", "glob of the directories to exclude")
	var optNoDefaultExcludes = flag.Bool("no-default-excludes", false, "don't apply default excludes")

	getopt.Aliases(
		"h", "help",
		"d", "depth",
		"f", "flag-file",
		"c", "check-file",
		"e", "check-regexp",
		"i", "check-inverse",
		"x", "exclude",
		"n", "no-default-excludes",
	)
	getopt.Parse()

	var optRootDir = "."
	if len(flag.Args()) > 0 {
		optRootDir = flag.Arg(0)
	}

	if *optHelp || len(*optFlagFile) == 0 {
		fmt.Println("Usage: ls-having -f name-or-glob [options] [root-dir]")
		fmt.Println("Options:")
		getopt.PrintDefaults()
		fmt.Println("References:")
		fmt.Println("  Glob syntax: https://github.com/gobwas/glob#example")
		fmt.Println("  Regexp syntax: https://pkg.go.dev/regexp/syntax")
		fmt.Println("  Home page: https://github.com/handy-common-utils/ls-having")
		if *optHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if !*optNoDefaultExcludes {
		optExcludes = append(optExcludes,
			".git",
			filepath.Join("**", ".git"),
			"node_modules",
			filepath.Join("**", "node_modules"),
		)
	}

	var dirs = lsh.LsHaving(&lsh.Options{
		Depth:        *optDepth,
		Excludes:     compileGlobs(optExcludes, filepath.Separator),
		FlagFile:     glob.MustCompile(filepath.Join(*optFlagFile), filepath.Separator),
		CheckFile:    *optCheckFile,
		CheckRegexp:  regexp.MustCompile(*optCheckRegexp),
		CheckInverse: *optCheckInverse,
	}, optRootDir)
	fmt.Println(dirs)
}

func compileGlobs(globStrings []string, separator rune) []glob.Glob {
	result := make([]glob.Glob, len(globStrings))
	for i, globString := range globStrings {
		result[i] = glob.MustCompile(globString, separator)
	}
	return result
}

type arrayFlag []string

func (i *arrayFlag) String() string {
	return strings.Join(*i, ",")
}
func (i *arrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
