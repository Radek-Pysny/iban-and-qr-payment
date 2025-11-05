package qrpay

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
)

const (
	spdSeparator = '*' // SID uses the same separator

	fieldSpdInitial = "SPD"
	fieldSidInitial = "SID"
	fieldVersion    = "version"

	versionSeparator  = "."
	decimalSeparator  = "."
	keyValueSeparator = ":"

	validValuePunctuation = "$%+-./,"
	extraValuePunctuation = " :"
)

type ModeFlags uint64

func (f ModeFlags) SkipPreprocessFailure() bool {
	return f&ModeFlagSkipPreprocessFailure != 0
}

func (f ModeFlags) SkipLengthCheck() bool {
	return f&ModeFlagSkipLengthCheck != 0
}

const (
	ModeFlagSkipPreprocessFailure ModeFlags = 1 << iota
	ModeFlagSkipLengthCheck

	ModeFlagsStandard       ModeFlags = 0
	ModeFlagsIgnoreFailures ModeFlags = ModeFlagSkipPreprocessFailure | ModeFlagSkipLengthCheck
)

func ShortPaymentDescriptorFromText(
	text string,
	flags ModeFlags,
	log *slog.Logger,
) (*ShortPaymentDescriptor, error) {
	header, major, minor, err := extractHeaderAndVersion(text)
	switch {
	case err != nil:
		return nil, err

	case header != fieldSpdInitial:
		return nil, NewParseError(fieldSpdInitial, header, "expected SID in the first field", nil)
	}

	result := ShortPaymentDescriptor{
		VersionMajor: major,
		VersionMinor: minor,
	}

	tagMap, err := prepareTagMap(result, tagKeySpd)
	if err != nil {
		return nil, NewInternalError("tag map preparation before parsing", err)
	}

	index := 0
	for field := range strings.FieldsFuncSeq(text, isSpdSeparator) {
		if index < 2 {
			index++
			continue // skip header and version fields
		}

		value, tag, err := parseKeyValuePair(field, tagMap, flags, log)
		if err != nil {
			return nil, err
		}

		rv := reflect.Indirect(reflect.ValueOf(&result))
		rf := rv.FieldByName(tag.structFieldName)
		rf.SetString(value)
	}

	if result.QrFaktura != "" {
		sid, err := ShortInvoiceDescriptorFromText(result.QrFaktura, flags, log)
		if err != nil {
			return nil, fmt.Errorf(`error parsing SID %q: %w`, result.QrFaktura, err)
		}

		result.SID = sid
	}

	return &result, nil
}

func ShortInvoiceDescriptorFromText(
	text string,
	flags ModeFlags,
	log *slog.Logger,
) (*ShortInvoiceDescriptor, error) {
	if text == "" {
		return nil, nil
	}

	header, major, minor, err := extractHeaderAndVersion(text)
	switch {
	case err != nil:
		return nil, err

	case header != fieldSidInitial:
		return nil, NewParseError(fieldSidInitial, header, "expected SID in the first field", nil)
	}

	result := ShortInvoiceDescriptor{
		VersionMajor: major,
		VersionMinor: minor,
	}

	tagMap, err := prepareTagMap(result, tagKeySid)
	if err != nil {
		return nil, NewInternalError("tag map preparation before parsing", err)
	}

	index := 0
	for field := range strings.FieldsFuncSeq(text, isSpdSeparator) {
		if index < 2 {
			index++
			continue // skip header and version fields
		}

		value, tag, err := parseKeyValuePair(field, tagMap, flags, log)
		if err != nil {
			if errors.Is(err, errSkip) {
				continue
			}

			return nil, err
		}

		rv := reflect.Indirect(reflect.ValueOf(&result))
		rf := rv.FieldByName(tag.structFieldName)
		rf.SetString(value)
	}

	return &result, nil
}

func IsShortPaymentDescriptor(text string) bool {
	header, _, _, err := extractHeaderAndVersion(text)

	return err == nil && header == fieldSpdInitial
}

func IsShortInvoiceDescriptor(text string) bool {
	header, _, _, err := extractHeaderAndVersion(text)

	return err == nil && header == fieldSidInitial
}

func extractHeaderAndVersion(text string) (header string, major uint8, minor uint8, err error) {
	index := 0
	for field := range strings.FieldsFuncSeq(text, isSpdSeparator) {
		index++

		switch index {
		case 1:
			header = field

		case 2:
			major, minor, err = parseVersion(field)

			return header, major, minor, err
		}
	}

	return "", 0, 0, errors.New("not enough fields for header and version extraction")
}

func parseKeyValuePair(
	text string,
	tagMap map[string]tagRecord,
	flags ModeFlags,
	log *slog.Logger,
) (string, tagRecord, error) {
	fieldKey, originalFieldValue, hasValue := strings.Cut(text, keyValueSeparator)
	if !hasValue {
		return "", tagRecord{}, NewParseError(fieldKey, "", "cannot detect value for field", nil)
	}

	tag, found := tagMap[fieldKey]
	if !found {
		if log != nil {
			log.Warn("not found in tag map", slog.String("key", fieldKey))
		}

		return "", tagRecord{}, errSkip
	}

	fieldValue, err := validateFieldValue(originalFieldValue)
	if err != nil {
		return "", tagRecord{}, NewParseError(fieldKey, originalFieldValue, "value preprocessing", err)
	}

	if tag.preprocessFn != nil {
		newFieldValue, err := tag.preprocessFn(fieldValue)
		if err != nil {
			if flags.SkipPreprocessFailure() {
				if log != nil {
					log.Warn(
						"skip preprocessing failure",
						slog.String("key", fieldKey),
						slog.String("value", fieldValue),
						slog.String("failure", err.Error()),
					)
				}
			} else {
				return "", tagRecord{}, NewParseError(fieldKey, fieldValue, "custom preprocessing", err)
			}
		}

		fieldValue = newFieldValue
	}

	{
		length := len(fieldValue)
		if length < int(tag.minCharCount) {
			err = fmt.Errorf("value is shorter than %d", tag.minCharCount)
		} else if length > int(tag.maxCharCount) {
			err = fmt.Errorf("value is longer than %d", tag.maxCharCount)
		}

		if err != nil {
			if flags.SkipLengthCheck() {
				if log != nil {
					log.Warn("skip length check failure",
						slog.String("key", fieldKey),
						slog.String("value", fieldValue),
						slog.Int("minLength", int(tag.minCharCount)),
						slog.Int("length", length),
						slog.Int("maxLength", int(tag.maxCharCount)),
					)
				}
			} else {
				return "", tagRecord{}, NewParseError(fieldKey, fieldValue, "length check", err)
			}
		}
	}

	return fieldValue, tag, nil
}

func isSpdSeparator(char rune) bool {
	return char == spdSeparator
}

func parseVersion(text string) (major uint8, minor uint8, err error) {
	v1, v2, found := strings.Cut(text, versionSeparator)
	if !found {
		return 0, 0, NewParseError(fieldVersion, text, `missing "." as separator`, nil)
	}

	var x uint64
	x, err = strconv.ParseUint(v1, 10, 8)
	if err != nil {
		return 0, 0, NewParseError(fieldVersion, text, `parse major number`, err)
	}
	major = uint8(x)

	x, err = strconv.ParseUint(v2, 10, 8)
	if err != nil {
		return 0, 0, NewParseError(fieldVersion, text, `parse minor number`, err)
	}
	minor = uint8(x)

	return major, minor, nil
}

func validateFieldValue(text string) (string, error) {
	for i, char := range text {
		switch {
		case char <= 0,
			char >= 255:
			return "", fmt.Errorf("not ASCII character '%c' at position %d", char, i)

		case char >= '0' && char <= '9',
			char >= 'A' && char <= 'Z',
			strings.ContainsRune(validValuePunctuation, char),
			strings.ContainsRune(extraValuePunctuation, char):
			break

		default:
			return "", fmt.Errorf("invalid ASCII character '%c' at position %d", char, i)
		}
	}

	return text, nil
}
