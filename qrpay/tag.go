package qrpay

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type TagKey string

const (
	tagKeySpd = TagKey("spd")
	tagKeySid = TagKey("sid")
)

func (k TagKey) String() string {
	return string(k)
}

const (
	tagSeparator             = ";"
	tagRequiredKey           = "req"
	tagOptionallyRequiredKey = "opt-req"
	tagOptionalKey           = "opt"
)

type TagRequired uint8

const (
	tagOptional           TagRequired = 0
	tagOptionallyRequired             = 1
	tagRequired                       = 2
)

type tagRecord struct {
	structFieldName    string
	originalTagContent string

	// info for SPD format
	fieldName    string
	required     TagRequired
	preprocessFn func(string) (string, error)
	minCharCount uint16
	maxCharCount uint16
	description  string
}

func prepareTagMap(x any, tagKey TagKey) (map[string]tagRecord, error) {
	typ := reflect.TypeOf(x)

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("ShortPaymentDescriptor is not a struct")
	}

	result := make(map[string]tagRecord, typ.NumField())

	for i := range typ.NumField() {
		field := typ.Field(i)

		tag, found := field.Tag.Lookup(tagKey.String())
		if !found {
			continue
		}

		record, err := parseTag(field.Name, tag)
		if err != nil {
			return nil, err
		}

		result[record.fieldName] = record
	}

	return result, nil
}

func parseTag(structFieldName string, text string) (tagRecord, error) {
	result := tagRecord{
		structFieldName:    structFieldName,
		originalTagContent: text,
	}

	// name;required
	name, rest, found := strings.Cut(text, tagSeparator)
	if !found {
		return result, errors.New("missing 1st separator (between name and required)")
	}
	if strings.ToUpper(name) != name {
		return result, fmt.Errorf("name %q must be uppercase", name)
	}

	result.fieldName = name

	var (
		preprocessFnAsText string
		requiredAsText     string
		minCharCountAsText string
		maxCharCountAsText string
		description        string
	)

	// required;preprocessFn
	requiredAsText, rest, found = strings.Cut(rest, tagSeparator)
	if !found {
		return result, errors.New(`missing 2nd separator (between required and preprocess function name): "` + text + `"`)
	}
	switch requiredAsText {
	case tagRequiredKey:
		result.required = tagRequired

	case tagOptionallyRequiredKey:
		result.required = tagOptionallyRequired

	case tagOptionalKey:
		result.required = tagOptional

	default:
		return result, fmt.Errorf(`2nd field of tag is not %q nor %q: %q`, tagRequired, tagOptional, text)
	}

	// preprocessFn;minCharCount
	preprocessFnAsText, rest, found = strings.Cut(rest, tagSeparator)
	if !found {
		return result, errors.New("missing 3rd separator (between preprocess function name and min char count)")
	}

	result.preprocessFn, found = preprocessFn[preprocessFnAsText]
	if !found {
		return result, fmt.Errorf("3rd field of tag contains unexpected preprocess function name %q", preprocessFnAsText)
	}

	// minCharCount;maxCharCount
	minCharCountAsText, rest, found = strings.Cut(rest, tagSeparator)
	if !found {
		return result, errors.New(`missing 4th separator (between min and max char count): "` + text + `"`)
	}

	minCharCount, err := strconv.ParseUint(minCharCountAsText, 10, 16)
	if err != nil {
		return result, fmt.Errorf("4th field is not a number: %w", err)
	}

	result.minCharCount = uint16(minCharCount)

	// maxCharCount;description
	maxCharCountAsText, description, found = strings.Cut(rest, tagSeparator)
	if !found {
		return result, errors.New(`missing 5th separator (between max char count and description): "` + text + `"`)
	}

	maxCharCount, err := strconv.ParseUint(maxCharCountAsText, 10, 16)
	if err != nil {
		return result, fmt.Errorf("5th field is not a number: %w", err)
	}
	if maxCharCount < minCharCount {
		return result, fmt.Errorf(
			"max char count %d < min char count %d",
			maxCharCount,
			minCharCount,
		)
	}

	result.maxCharCount = uint16(maxCharCount)
	result.description = description

	return result, nil
}
