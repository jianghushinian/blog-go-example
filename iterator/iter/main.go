// NOTE: iter 包

package main

import (
	"iter"
	_ "maps"
	_ "slices"
)

func iter1[Slice ~[]E, E any](s Slice) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func iter2[Slice ~[]E, E any](s Slice) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i, v := range s {
			if !yield(i, v) {
				return
			}
		}
	}
}
