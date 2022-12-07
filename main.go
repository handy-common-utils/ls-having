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
	setupFlags()                       // this function can't be called more than one time globally
	doMain(println, printUsageAndExit) // this function is also called in every test case
}

const DEFAULT_DEPTH = 5

var optHelp *bool
var optDepth *int
var optFlagFiles arrayFlag
var optCheckFile *string
var optCheckRegexp *string
var optCheckInverse *bool
var optExcludes arrayFlag
var optNoDefaultExcludes *bool
var optOnlySubdirectories *bool

func setupFlags() {
	optHelp = flag.Bool("help", false, "show help information")
	optDepth = flag.Int("depth", DEFAULT_DEPTH, "how deep to look into subdirectories, 0 means only look at root directory, -1 means no limit")
	flag.Var(&optFlagFiles, "flag-file", "name or `glob` of the flag file, this option can appear multiple times")
	optCheckFile = flag.String("check-file", "", "`name` of the additional file to check")
	optCheckRegexp = flag.String("check-regexp", ".*", "regular `expression` for testing the content of the check file")
	optCheckInverse = flag.Bool("check-inverse", false, "regard regular expression not matching as positive")
	flag.Var(&optExcludes, "exclude", "`glob` of the directories to exclude, this option can appear multiple times")
	optNoDefaultExcludes = flag.Bool("no-default-excludes", false, "don't apply default excludes")
	optOnlySubdirectories = flag.Bool("subdirectories-only", false, "don't return root directory even if it meets conditions")

	getopt.Aliases(
		"h", "help",
		"d", "depth",
		"f", "flag-file",
		"c", "check-file",
		"e", "check-regexp",
		"i", "check-inverse",
		"x", "exclude",
		"n", "no-default-excludes",
		"s", "subdirectories-only",
	)
	flag.Usage = func() {
		// do nothing, just to avoid getopt to show usage after warning/error info
	}
}

func parseFlags() {
	// It seems that getopt.Parse() does not reset flags in case it is called more than one time
	*optHelp = false
	*optDepth = DEFAULT_DEPTH
	optFlagFiles = nil
	*optCheckFile = ""
	*optCheckRegexp = ""
	*optCheckInverse = false
	optExcludes = nil
	*optNoDefaultExcludes = false
	*optOnlySubdirectories = false

	getopt.Parse()
}

func doMain(println func(text string), printUsageAndExit func(text string)) {
	parseFlags()

	if *optHelp {
		printUsageAndExit("")
		return
	}

	var optRootDir = getopt.CommandLine.Arg(0)
	if optRootDir == "" {
		optRootDir = "."
	}

	if len(optFlagFiles) == 0 {
		printUsageAndExit("flag file has not been specified")
		return
	}

	if !*optNoDefaultExcludes {
		optExcludes = append(optExcludes,
			".git",
			filepath.Join("**", ".git"),
			"node_modules",
			filepath.Join("**", "node_modules"),
		)
	}

	var options = lsh.Options{
		Depth:        *optDepth,
		Excludes:     compileGlobs(optExcludes, filepath.Separator),
		ExcludeRoot:  *optOnlySubdirectories,
		FlagFiles:    compileGlobs(optFlagFiles, filepath.Separator),
		CheckFile:    *optCheckFile,
		CheckRegexp:  regexp.MustCompile(*optCheckRegexp),
		CheckInverse: *optCheckInverse,
	}
	var dirs = lsh.LsHaving(&options, optRootDir)
	if len(dirs) > 0 {
		println(strings.Join(dirs, "\n"))
	}
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
	return "xyz" // strings.Join(*i, ",")
}
func (i *arrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func println(text string) {
	fmt.Println(text)
}

func printUsageAndExit(errorString string) {
	var writer = os.Stdout
	var exitCode = 0
	if len(errorString) > 0 {
		writer = os.Stderr
		exitCode = 1
		fmt.Fprintln(writer, "Error: "+errorString)
	}
	fmt.Println("Usage: ls-having -f name-or-glob [options] [root-dir]")
	fmt.Println("Options:")
	getopt.CommandLine.SetOutput(os.Stdout)
	getopt.PrintDefaults()
	fmt.Println("References:")
	fmt.Println("  Glob syntax: https://github.com/gobwas/glob#example")
	fmt.Println("  Regexp syntax: https://pkg.go.dev/regexp/syntax")
	fmt.Println("  Home page: https://github.com/handy-common-utils/ls-having")
	os.Exit(exitCode)
}
