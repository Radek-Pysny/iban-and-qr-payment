package stringutils

import (
	"iter"
)

func SeparateGroups(input string, groupSize int, separator string) string {
	sep := ""
	result := ""
	for group := range StringChunks(input, groupSize) {
		result += sep + group
		sep = separator
	}

	return result
}

func StringChunks(input string, chunkSize int) iter.Seq[string] {
	if chunkSize < 1 {
		return func(yield func(string) bool) {
			yield(input)
		}
	}

	return func(yield func(string) bool) {
		start := 0
		length := len(input)

		for start < length {
			end := start + chunkSize
			if end > length {
				end = length
			}

			if !yield(input[start:end]) {
				return
			}

			start = end
		}
	}
}
