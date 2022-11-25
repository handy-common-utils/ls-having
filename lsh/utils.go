package lsh

import (
	"github.com/gobwas/glob"
)

func anyGlobMatch(globs []glob.Glob, text string) bool {
	for _, glob := range globs {
		if glob.Match(text) {
			return true
		}
	}
	return false
}
