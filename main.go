package main

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/gobwas/glob"
	"github.com/handy-common-utils/ls-having/lsh"
)

func main() {
	var dirs = lsh.LsHaving(&lsh.Options{
		Depth: 10,
		Excludes: compileGlobs([]string{
			".git",
			filepath.Join("**", ".git"),
			"node_modules",
			filepath.Join("**", "node_modules"),
		}, filepath.Separator),
		FlagFile:     glob.MustCompile(filepath.Join("package.*"), filepath.Separator),
		CheckFile:    "package.json",
		CheckRegexp:  regexp.MustCompile("(?m)^ *\"mocha"),
		CheckInverse: false,
	}, ".")
	fmt.Println(dirs)
}

func compileGlobs(globStrings []string, separator rune) []glob.Glob {
	result := make([]glob.Glob, len(globStrings))
	for i, globString := range globStrings {
		result[i] = glob.MustCompile(globString, separator)
	}
	return result
}
