package lsh

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/gobwas/glob"
)

type Options struct {
	Depth        int
	Excludes     []glob.Glob
	ExcludeRoot  bool
	FlagFiles    []glob.Glob
	CheckFile    string // "" means it is not used
	CheckRegexp  *regexp.Regexp
	CheckInverse bool
}

func LsHaving(options *Options, dir string) []string {
	var found []string = make([]string, 0, 100)
	doLsHaving(options, &found, dir, 0, nil) // root dir has depth 0
	sort.Strings(found)
	return found
}

type dirEntryEx struct {
	Path  string
	Depth int
	Entry fs.DirEntry
}

// Read all entries under the specified directory.
// The "depth" parameter is the depth of the directory specified by "dir" parameter.
func readEntries(dir string, depth int) *[]dirEntryEx {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err) // panic if can't list the content in the directory
	}
	entriesEx := make([]dirEntryEx, 0, len(entries))
	for _, entry := range entries {
		entriesEx = append(entriesEx, dirEntryEx{filepath.Join(dir, entry.Name()), depth + 1, entry})
	}
	return &entriesEx
}

func doLsHaving(options *Options, found *[]string, dir string, depth int, entriesInDir *[]dirEntryEx) {
	if entriesInDir == nil { // this must be the root dir
		if depth != 0 {
			log.Fatal("Internal error: entriesInDir == nil but depth != 0")
		}
		rootDirInfo, err := os.Stat(dir)
		if err != nil {
			log.Fatal(err) // panic if can't read the root directory
		}
		rootDirEntryEx := dirEntryEx{dir, 0, fs.FileInfoToDirEntry(rootDirInfo)}
		entriesInDir = readEntries(dir, 0)

		if shouldCheck(options, &rootDirEntryEx) {
			if match(options, &rootDirEntryEx, entriesInDir) {
				*found = append(*found, rootDirEntryEx.Path)
			}
		}
	}

	for _, entry := range *entriesInDir {
		if shouldCheck(options, &entry) {
			entriesInSubDir := readEntries(entry.Path, entry.Depth)
			if match(options, &entry, entriesInSubDir) {
				*found = append(*found, entry.Path)
			}
			doLsHaving(options, found, entry.Path, entry.Depth, entriesInSubDir)
		}
	}
}

func shouldCheck(options *Options, dir *dirEntryEx) bool {
	if !dir.Entry.IsDir() {
		return false
	}
	if options.Depth >= 0 && options.Depth < dir.Depth {
		return false
	}
	if options.Excludes != nil && anyGlobMatch(options.Excludes, dir.Path) {
		return false
	}
	return true
}

func match(options *Options, dir *dirEntryEx, entries *[]dirEntryEx) bool {
	if options.ExcludeRoot && dir.Depth == 0 {
		return false
	}

	foundFlagFile := false
	for _, entry := range *entries {
		if anyGlobMatch(options.FlagFiles, entry.Entry.Name()) {
			foundFlagFile = true
			break
		}
	}
	var checkFileMismatch bool
	if options.CheckFile == "" {
		checkFileMismatch = false
	} else {
		checkFilePath := filepath.Join(dir.Path, options.CheckFile)
		checkFileContent, err := os.ReadFile(checkFilePath)
		if err != nil {
			checkFileMismatch = true
		} else {
			checkFileMismatch = !options.CheckRegexp.Match(checkFileContent)
			checkFileMismatch = checkFileMismatch != options.CheckInverse
		}
	}
	return foundFlagFile && !checkFileMismatch
}
