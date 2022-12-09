// This package exposes the functionality of ls-having.
// For documentation of ls-having,
// see its home page: https://github.com/handy-common-utils/ls-having
package lsh

import (
	"io/fs"
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

	// To return immedately when any error (such like non-existing directory or no access permission) happens
	PanicOnError bool
}

// Find directories matching conditions.
//
// The array of paths is returned as the first value,
// it could be an empty array but would never be nil,
// and the array is sorted in ascend order.
// If there is no error, the second value returned would be nil.
// If there is any error, the array of error messages is returned as the second value.
//
// If the options tells this function to panic on error,
// the function would return immediately once there's an error.
// In such case the first returned value could contain some paths,
// and the second returned value would contain the error message.
//
// If the options tells this function to not panic on error,
// the function would record the error but continue working in case an error happens.
// In such case the first returned value would contain all paths found,
// and the second returned value would contain the error messages.
func LsHaving(options *Options, rootDir string) (found []string, errors []string) {
	found = make([]string, 0, 100)
	errors = make([]string, 0, 10)

	doLsHaving(options, &found, &errors, rootDir, 0, nil) // root dir has depth 0

	sort.Strings(found)
	if len(errors) == 0 {
		errors = nil
	}
	return
}

type dirEntryEx struct {
	Path  string
	Depth int
	Entry fs.DirEntry
}

// Read all entries under the specified directory.
// The "depth" parameter is the depth of the directory specified by "dir" parameter.
//
// In case any error happens, the returned values would have an empty array and the error
func readEntries(dir string, depth int) (*[]dirEntryEx, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return &[]dirEntryEx{}, err
	}
	entriesEx := make([]dirEntryEx, 0, len(entries))
	for _, entry := range entries {
		entriesEx = append(entriesEx, dirEntryEx{filepath.Join(dir, entry.Name()), depth + 1, entry})
	}
	return &entriesEx, nil
}

func doLsHaving(options *Options, found *[]string, errors *[]string, dir string, depth int, entriesInDir *[]dirEntryEx) {
	if entriesInDir == nil { // this must be the root dir
		rootDirInfo, err := os.Stat(dir)
		if err != nil {
			*errors = append(*errors, err.Error())
			return // If root dir cannot be read, there's no point to continue even if options.PanicOnError == false
		}
		rootDirEntryEx := dirEntryEx{dir, 0, fs.FileInfoToDirEntry(rootDirInfo)}
		entriesInDir, err = readEntries(dir, 0)
		if err != nil {
			*errors = append(*errors, err.Error())
			if options.PanicOnError {
				return
			}
		}

		if shouldCheck(options, &rootDirEntryEx) {
			if match(options, &rootDirEntryEx, entriesInDir) {
				*found = append(*found, rootDirEntryEx.Path)
			}
		}
	}

	for _, entry := range *entriesInDir {
		if shouldCheck(options, &entry) {
			entriesInSubDir, err := readEntries(entry.Path, entry.Depth)
			if err != nil {
				*errors = append(*errors, err.Error())
				if options.PanicOnError {
					return
				}
			}
			if match(options, &entry, entriesInSubDir) {
				*found = append(*found, entry.Path)
			}
			doLsHaving(options, found, errors, entry.Path, entry.Depth, entriesInSubDir)
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
		checkFileDirInfo, err := os.Stat(checkFilePath)
		if err != nil {
			// can't find or cannot read check file/dir
			checkFileMismatch = !options.CheckInverse
		} else {
			if checkFileDirInfo.IsDir() { // it is a directory
				// ".*" is the default matching all expression
				checkFileMismatch = options.CheckRegexp.String() == ".*" == options.CheckInverse
			} else { // it is a file
				checkFileContent, err := os.ReadFile(checkFilePath)
				if err != nil {
					checkFileMismatch = true
				} else {
					checkFileMismatch = !options.CheckRegexp.Match(checkFileContent)
					checkFileMismatch = checkFileMismatch != options.CheckInverse
				}
			}
		}
	}
	return foundFlagFile && !checkFileMismatch
}
