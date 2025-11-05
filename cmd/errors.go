package cmd

import (
	"errors"
	"fmt"
)

var (
	errEmptyInput = errors.New("empty input")

	errInternal   = errors.New("internal error")
	errReadLn     = fmt.Errorf("%w: unable to read line", errInternal)
	errTypeAssert = fmt.Errorf("%w: type assertion", errInternal)
)
