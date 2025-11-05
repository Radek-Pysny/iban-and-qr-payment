package qrpay

import (
	"errors"
	"strings"
)

// errSkip is sentinel used to command continue to the for loop in a caller.
var errSkip = errors.New("skip")

type ErrorKind uint8

const (
	ErrorKindInvalidArgument ErrorKind = iota
	ErrorKindInternal
)

func (k ErrorKind) String() string {
	switch k {
	case ErrorKindInvalidArgument:
		return "INVALID_ARGUMENT_ERROR"

	case ErrorKindInternal:
		return "INTERNAL_ERROR"

	default:
		return "UNKNOWN_ERROR"
	}
}

type baseError struct {
	errorKind     ErrorKind
	originalError error
	customError   string
}

func (e *baseError) Error() string {
	return showError(e.errorKind, e.customError, e.originalError)
}

func (e *baseError) Equal(other *baseError) bool {
	switch {
	case e.errorKind != other.errorKind,
		e.customError != other.customError:
		return false

	case e.originalError == nil && other.originalError == nil:
		return true

	case e.originalError == nil,
		other.originalError == nil,
		e.originalError.Error() != other.originalError.Error():
		return false

	default:
		return true
	}
}

type ParseKeyError struct {
	key           string
	content       string
	customError   string
	originalError error
}

func (e *ParseKeyError) String() string {
	return `key ` + e.key + ` = "` + e.content + `"` + `: ` + showErrorDescription(e.customError, e.originalError)
}

type ShortPaymentDescriptorParseError struct {
	baseError
	keyErrors []ParseKeyError
}

func (e *ShortPaymentDescriptorParseError) Error() string {
	result := showError(e.errorKind, e.customError, e.originalError)
	for i := range e.keyErrors {
		keyError := &e.keyErrors[i]
		result += "\n  - " + keyError.String()
	}
	if len(e.keyErrors) > 0 {
		result += "\n"
	}

	return result
}

func (e *ShortPaymentDescriptorParseError) Equal(other *ShortPaymentDescriptorParseError) bool {
	if !e.baseError.Equal(&other.baseError) {
		return false
	}

	for i, keyError := range e.keyErrors {
		otherKeyError := &other.keyErrors[i]

		switch {
		case keyError.key != otherKeyError.key,
			keyError.content != otherKeyError.content,
			keyError.customError != otherKeyError.customError:
			return false

		case keyError.originalError == nil && otherKeyError.originalError == nil:
			return true

		case keyError.originalError == nil,
			otherKeyError.originalError == nil,
			keyError.originalError.Error() != otherKeyError.originalError.Error():
			return false
		}
	}

	return true
}

func NewInternalError(
	customError string,
	originalError error,
) error {
	return &baseError{
		errorKind:     ErrorKindInternal,
		originalError: originalError,
		customError:   customError,
	}
}

func NewParseError(
	key string,
	content string,
	customError string,
	originalError error,
) *ShortPaymentDescriptorParseError {
	return &ShortPaymentDescriptorParseError{
		baseError: baseError{
			errorKind:   ErrorKindInvalidArgument,
			customError: "failed to parse SPD or SID format",
		},
		keyErrors: []ParseKeyError{
			ParseKeyError{
				key:           key,
				content:       content,
				customError:   customError,
				originalError: originalError,
			},
		},
	}
}

func (e *ShortPaymentDescriptorParseError) Append(
	key string,
	content string,
	customError string,
	originalError error,
) *ShortPaymentDescriptorParseError {
	e.keyErrors = append(e.keyErrors, ParseKeyError{
		key:           key,
		content:       content,
		customError:   customError,
		originalError: originalError,
	})

	return e
}

func showError(errorKind ErrorKind, customError string, originalError error) string {
	customError = strings.TrimSpace(customError)

	if customError == "" && originalError == nil {
		return "no error"
	}

	return "[" + errorKind.String() + "] " + showErrorDescription(customError, originalError)
}

func showErrorDescription(customError string, originalError error) string {
	result := ""
	sep := ""

	if customError != "" {
		result += customError
		sep = ": "
	}

	if originalError != nil {
		result += sep + originalError.Error()
	}

	return result
}
