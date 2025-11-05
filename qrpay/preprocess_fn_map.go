package qrpay

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var preprocessFn = map[string]func(string) (string, error){
	"-":  nil,
	"no": nil,

	// replacement of "%2A" → "*"
	"ast": func(originalText string) (string, error) {
		if strings.IndexByte(originalText, '*') >= 0 {
			return originalText, fmt.Errorf("ast preprocess fn: found '*' in %q", originalText)
		}

		newText := strings.ReplaceAll(originalText, "%2A", "*")

		if !strings.HasPrefix(newText, "SID*") {
			return originalText, fmt.Errorf("ast preprocess fn: missing 'SID%%2A' prefix in %q", originalText)
		}

		return newText, nil
	},

	// date in standard format YYYYMMDD
	"date": func(originalText string) (string, error) {
		if _, err := time.Parse("20060102", originalText); err != nil {
			return originalText, fmt.Errorf("date preprocess fn: invalid date %q: %v", originalText, err)
		}

		return originalText, nil
	},

	// Digits only
	"dig": func(originalText string) (string, error) {
		if !reDig.MatchString(originalText) {
			return originalText, fmt.Errorf("dig preprocess fn: invalid format of %q", originalText)
		}

		return originalText, nil
	},

	// Digit Boolean flag
	"dig-bool": func(originalText string) (string, error) {
		if originalText != "1" && originalText != "0" {
			return originalText, fmt.Errorf("dig-bool preprocess fn: invalid format of %q", originalText)
		}

		return originalText, nil
	},

	// number in fp format
	"num-fp": func(originalText string) (string, error) {
		if !reNumFP.MatchString(originalText) {
			return originalText, fmt.Errorf("num-fp preprocess fn: invalid format of %q", originalText)
		}

		return originalText, nil
	},

	// abc (only letters)
	"abc": func(originalText string) (string, error) {
		if !reLetter.MatchString(originalText) {
			return originalText, fmt.Errorf("abc preprocess fn: invalid format of %q", originalText)
		}

		return originalText, nil
	},

	// hex (hexadecimal digits)
	"hex": func(originalText string) (string, error) {
		if !reHex.MatchString(originalText) {
			return originalText, fmt.Errorf("abc preprocess fn: invalid format of %q", originalText)
		}

		return originalText, nil
	},
}

var reDig = regexp.MustCompile(`^[0-9]*$`)

var reHex = regexp.MustCompile(`^[0-9A-F]*$`)

var reLetter = regexp.MustCompile(`^[A-Z]*$`)

var reNumFP = regexp.MustCompile(`^[0-9]+(\.[0-9][0-9]?)?$`)
