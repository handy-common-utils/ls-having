package lsh

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gobwas/glob"
)

type Options struct {
	Depth            int
	Exclude          glob.Glob
	NoDefaultExclude bool
	FlagFile         string
	CheckFile        string
	CheckRegexp      regexp.Regexp
	CheckInverse     bool
}

func LsHaving(options *Options, dir string) []string {
	var found []string = make([]string, 0, 100)
	doLsHaving(options, &found, dir, nil)
	return found
}

type dirEntryEx struct {
	Path  string
	Entry fs.DirEntry
}

func readEntries(dir string) *[]dirEntryEx {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	entriesEx := make([]dirEntryEx, 0, len(entries))
	for _, entry := range entries {
		entriesEx = append(entriesEx, dirEntryEx{filepath.Join(dir, entry.Name()), entry})
	}
	return &entriesEx
}

func doLsHaving(options *Options, found *[]string, dir string, entriesInDir *[]dirEntryEx) {
	if entriesInDir == nil {
		entriesInDir = readEntries(dir)
	}
	for _, entry := range *entriesInDir {
		if shouldCheck(options, &entry) {
			entriesInSubDir := readEntries(entry.Path)
			if match(options, &entry, entriesInSubDir) {
				*found = append(*found, entry.Path)
			}
			doLsHaving(options, found, entry.Path, entriesInSubDir)
		}
	}
}

func shouldCheck(options *Options, dir *dirEntryEx) bool {
	if !dir.Entry.IsDir() {
		return false
	}
	return true
}

func match(options *Options, dir *dirEntryEx, entries *[]dirEntryEx) bool {
	for _, entry := range *entries {
		if entry.Entry.Name() == "serverless.yml" {
			return true
		}
	}
	return false
}
