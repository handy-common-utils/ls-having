// This package ("main" package) is not supposed to be imported by another program.
// To use the functionality of ls-having, the "lsh" package
// ("github.com/handy-common-utils/ls-having/lsh" package) can be imported.
// See the documentation of "lsh" package for details.
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
	doMain(printToStdout, handleError) // this function is also called in every test case
}

const OPT_ERROR_PANIC = "panic"
const OPT_ERROR_IGNORE = "ignore"
const OPT_ERROR_PRINT = "print"

const DEFAULT_DEPTH = 5
const DEFAULT_CHECK_REGEXP = ".*"
const DEFAULT_ERROR = OPT_ERROR_IGNORE

const DEFAULT_EXIT_CODE_WHEN_ERROR = 1

var optHelp *bool
var optDepth *int
var optFlagFiles arrayFlag
var optCheckFile *string
var optCheckRegexp *string
var optCheckInverse *bool
var optExcludes arrayFlag
var optNoDefaultExcludes *bool
var optOnlySubdirectories *bool
var optPrint0 *bool
var optError *string

func setupFlags() {
	optHelp = flag.Bool("help", false, "show help information")
	optDepth = flag.Int("depth", DEFAULT_DEPTH, "how deep to look into subdirectories, 0 means only look at root directory, -1 means no limit")
	flag.Var(&optFlagFiles, "flag-file", "name or `glob` of the flag file, this option can appear multiple times")
	optCheckFile = flag.String("check-file", "", "`name` of the additional file to check")
	optCheckRegexp = flag.String("check-regexp", DEFAULT_CHECK_REGEXP, "regular `expression` for testing the content of the check file")
	optCheckInverse = flag.Bool("check-inverse", false, "regard regular expression not matching as positive")
	flag.Var(&optExcludes, "exclude", "`glob` of the directories to exclude, this option can appear multiple times")
	optNoDefaultExcludes = flag.Bool("no-default-excludes", false, "don't apply default excludes")
	optOnlySubdirectories = flag.Bool("subdirectories-only", false, "don't return root directory even if it meets conditions")
	optPrint0 = flag.Bool("print0", false, "separate paths in the output with null characters (instead of newline characters)")
	optError = flag.String("error", DEFAULT_ERROR, "how to handle errors such like non-existing directory, no access permission, etc. (ignore|panic|print)")

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
		"0", "print0",
		"r", "error",
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
	*optCheckRegexp = DEFAULT_CHECK_REGEXP
	*optCheckInverse = false
	optExcludes = nil
	*optNoDefaultExcludes = false
	*optOnlySubdirectories = false
	*optPrint0 = false

	getopt.Parse()
}

func doMain(printOutput func(text string), handleError func(errors []string, printUsage bool, exitCode int)) {
	parseFlags()

	if *optHelp {
		handleError(nil, true, 0)
		return
	}

	var optRootDir = getopt.CommandLine.Arg(0)
	if optRootDir == "" {
		optRootDir = "."
	}

	if len(optFlagFiles) == 0 {
		if len(*optCheckFile) == 0 {
			handleError([]string{"flag file or check file must be specified"}, true, DEFAULT_EXIT_CODE_WHEN_ERROR)
			return
		} else {
			// assuming the check file is also the flag file
			optFlagFiles = append(optFlagFiles, *optCheckFile)
		}
	}

	if !*optNoDefaultExcludes {
		optExcludes = append(optExcludes,
			".git",
			filepath.Join("**", ".git"),
			"node_modules",
			filepath.Join("**", "node_modules"),
			"testdata",
			filepath.Join("**", "testdata"),
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
		PanicOnError: *optError == OPT_ERROR_PANIC,
	}
	var dirs, errors = lsh.LsHaving(&options, optRootDir)
	if errors != nil {
		switch *optError {
		case OPT_ERROR_PANIC:
			handleError(errors, false, DEFAULT_EXIT_CODE_WHEN_ERROR)
			return
		case OPT_ERROR_PRINT:
			handleError(errors, false, 0)
			// continue to print out results
		default:
			// do nothing
			// continue to print out results
		}
	}
	if len(dirs) > 0 {
		separator := "\n"
		if *optPrint0 {
			separator = string([]byte{0})
		}
		printOutput(strings.Join(dirs, separator))
		printOutput(separator)
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
	return strings.Join(*i, ",")
}
func (i *arrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func printToStdout(text string) {
	fmt.Print(text)
}

// Handle error situation.
// Depending on the parameters passed in, this function could
//   - print out error messages to stderr
//   - print out usage/help info to stdout
//   - exit the program with a non-zero exit code
//
// Parameters:
//   - errors: nil or an array of error messages which would be printed to stderr
//   - printUsage: true if usage/help info should be printed to stdout
//   - exitCode: non-zero exit code if os.Exit should be called, or zero if this function should return normally
func handleError(errors []string, printUsage bool, exitCode int) {
	var writer = os.Stdout
	if len(errors) > 0 {
		writer = os.Stderr
		for _, errorString := range errors {
			fmt.Fprintln(writer, "Error: "+errorString)
		}
	}
	if printUsage {
		fmt.Println("Usage: ls-having -f name-or-glob [options] [root-dir]")
		fmt.Println("Options:")
		getopt.CommandLine.SetOutput(os.Stdout)
		getopt.PrintDefaults()
		fmt.Println("References:")
		fmt.Println("  Glob syntax: https://github.com/gobwas/glob#example")
		fmt.Println("  Regexp syntax: https://pkg.go.dev/regexp/syntax")
		fmt.Println("  Home page: https://github.com/handy-common-utils/ls-having")
	}
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
