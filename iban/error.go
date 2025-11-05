package iban

import (
	"errors"
)

var (
	ErrWrongLength      = errors.New("wrong length")
	ErrWrongPrefix      = errors.New("wrong prefix")
	ErrExpectedDigits   = errors.New("expected only digits")
	ErrWrongCheckDigits = errors.New("wrong check digits")
	ErrWrongCharacter   = errors.New("wrong character")
)
