package lsh

import (
	"github.com/gobwas/glob"
)

type emptyStruct struct{}

var emptyVar emptyStruct

func anyGlobMatch(globs []glob.Glob, text string) bool {
	for _, glob := range globs {
		if glob.Match(text) {
			return true
		}
	}
	return false
}

func allMatchingGlobs(globs []glob.Glob, text string) map[int]emptyStruct {
	set := make(map[int]emptyStruct)
	for i, glob := range globs {
		if glob.Match(text) {
			set[i] = emptyVar
		}
	}
	return set
}
