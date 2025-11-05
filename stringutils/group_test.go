package stringutils_test

import (
	"testing"

	"github.com/Radek-Pysny/iban-and-qr-payment/stringutils"
	"github.com/stretchr/testify/assert"
)

func Test_SeparateGroup(t *testing.T) {
	testCases := []struct {
		title     string
		input     string
		chunkSize int
		separator string
		expected  string
	}{
		{
			title:     "empty-input-negative-chunk-size",
			input:     "",
			chunkSize: -10,
			separator: " ",
			expected:  "",
		},
		{
			title:     "empty-input-zero-chunk-size",
			input:     "",
			chunkSize: 0,
			separator: " ",
			expected:  "",
		},
		{
			title:     "empty-input-positive-chunk-size",
			input:     "",
			chunkSize: 0,
			separator: " ",
			expected:  "",
		},
		{
			title:     "empty-size4-space",
			input:     "",
			chunkSize: 4,
			separator: " ",
			expected:  "",
		},
		{
			title:     "abc-size0-space",
			input:     "abc",
			chunkSize: 0,
			separator: " ",
			expected:  "abc",
		},
		{
			title:     "abc-negative-size-space",
			input:     "abc",
			chunkSize: -2,
			separator: " ",
			expected:  "abc",
		},
		{
			title:     "abc-size1-space",
			input:     "abc",
			chunkSize: 1,
			separator: " ",
			expected:  "a b c",
		},
		{
			title:     "abc-size2-space",
			input:     "abc",
			chunkSize: 2,
			separator: " ",
			expected:  "ab c",
		},
		{
			title:     "abc-size3-space",
			input:     "abc",
			chunkSize: 3,
			separator: " ",
			expected:  "abc",
		},
		{
			title:     "abc-size4-space",
			input:     "abc",
			chunkSize: 4,
			separator: " ",
			expected:  "abc",
		},
		{
			title:     "hex-size4-space",
			input:     "0123456789ABCDEF",
			chunkSize: 4,
			separator: " ",
			expected:  "0123 4567 89AB CDEF",
		},
		{
			title:     "hex-size5-space",
			input:     "0123456789ABCDEF",
			chunkSize: 5,
			separator: " ",
			expected:  "01234 56789 ABCDE F",
		},
		{
			title:     "hex-size5-space",
			input:     "0123456789ABCDEF",
			chunkSize: 6,
			separator: " ",
			expected:  "012345 6789AB CDEF",
		},
		{
			title:     "hex-size5-space",
			input:     "0123456789ABCDEF",
			chunkSize: 7,
			separator: " ",
			expected:  "0123456 789ABCD EF",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			returned := stringutils.SeparateGroups(tc.input, tc.chunkSize, tc.separator)

			assert.Equal(t, tc.expected, returned)
		})
	}
}
