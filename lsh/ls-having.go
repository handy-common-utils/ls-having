// This package exposes the functionality of ls-having.
// For documentation of ls-having,
// see its home page: https://github.com/handy-common-utils/ls-having
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

// For more details regarding these options,
// see the home page of ls-having: https://github.com/handy-common-utils/ls-having
type Options struct {
	// Maximum depth to look into subdirectories.
	// The root directory has depth 0.
	// Negative value in this field means no limitation on depth.
	Depth int

	// Directories matching any of these patterns won't be looked into
	Excludes []glob.Glob

	// Exclude root directory in the result to be returned
	ExcludeRoot bool

	// Only directories having at least one file matching any of these patterns could be returned
	FlagFiles []glob.Glob

	// Additional file that its content would be checked. Use empty string to skip this checking.
	CheckFile string

	// Regular expression used for checking the content of CheckFile
	CheckRegexp *regexp.Regexp

	// Regard not matching as positive when using CheckRegexp to check the content of CheckFile
	CheckInverse bool
}

// Find directories matching conditions.
// The array of paths returned is sorted in ascend order.
func LsHaving(options *Options, rootDir string) []string {
	var found []string = make([]string, 0, 100)
	doLsHaving(options, &found, rootDir, 0, nil) // root dir has depth 0
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
