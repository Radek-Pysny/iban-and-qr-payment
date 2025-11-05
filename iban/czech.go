package iban

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Radek-Pysny/iban-and-qr-payment/stringutils"
)

const (
	validLength             = 24
	maxValueOfBankCode      = 9_999         // 4 digits
	maxValueOfAccountPrefix = 999_999       // 6 digits
	maxValueOfAccountBase   = 9_999_999_999 // 10 digits
)

var (
	ErrCzNoSepBetweenAccountNumberAndBankCode = errors.New("missing '/' to separate account number and bank code")
	ErrCzParseAccountNumberPrefix             = errors.New("failed to parse account number prefix")
	ErrCzParseAccountNumberBase               = errors.New("failed to parse account number base part")
	ErrCzParseBankCode                        = errors.New("failed to parse bank code")
	ErrCzTooLargeAccountNumberPrefix          = errors.New("account number prefix is too large")
	ErrCzTooLargeAccountNumberBase            = errors.New("account number base is too large")
	ErrCzTooLargeBankCode                     = errors.New("bank code is too large")
	ErrCzZeroAccountNumberBase                = errors.New("account number base is zero")
	ErrCzZeroBankCode                         = errors.New("bank code is zero")
)

type CzechIBAN struct {
	accountNumberBase   uint32
	accountNumberPrefix uint32
	bankCode            uint16
	checkDigits         uint8
}

func NewCzechIBAN(
	accountNumberPrefix uint32,
	accountNumberBase uint32,
	bankCode uint16,
	checkDigits uint8,
) (CzechIBAN, error) {
	result := CzechIBAN{
		accountNumberBase:   accountNumberBase,
		accountNumberPrefix: accountNumberPrefix,
		bankCode:            bankCode,
		checkDigits:         checkDigits,
	}

	if err := checkAccountNumberPrefix(uint64(accountNumberPrefix)); err != nil {
		return result, err
	}

	if err := checkAccountNumberBase(uint64(accountNumberBase)); err != nil {
		return result, err
	}

	if err := checkBankCode(uint64(accountNumberBase)); err != nil {
		return result, err
	}

	// TODO: check validity...

	return result, nil
}

// CzechIbanFromText parse Czech IBAN from either machine-readable form (aka electronic form; without spaces)
// or human-readable form (aka written form; with groups of 4 characters each).
// So all the following inputs are accepted:
// - CZ6508000000192000145399
// - CZ65 0800 0000 1920 0014 5399
func CzechIbanFromText(text string) (CzechIBAN, error) {
	text = strings.ReplaceAll(text, " ", "")

	result := CzechIBAN{}

	if len(text) != validLength {
		return result, ErrWrongLength
	}

	if !strings.HasPrefix(text, "CZ") {
		return result, ErrWrongPrefix
	}

	numbers := text[2:]
	if strings.ContainsFunc(numbers, func(r rune) bool { return r < '0' || r > '9' }) {
		return result, fmt.Errorf("%w: %q", ErrExpectedDigits, numbers)
	}

	checkDigits := text[2:4]
	bankCode := text[4:8]
	accountNumberPrefix := text[8:14]
	accountNumberBase := text[14:]

	prefix, base, code, err := parseEachPart(accountNumberPrefix, accountNumberBase, bankCode)
	if err != nil {
		return result, err
	}

	x, _ := strconv.ParseUint(checkDigits, 10, 8)

	result.checkDigits = uint8(x)
	result.accountNumberPrefix = prefix
	result.accountNumberBase = base
	result.bankCode = code

	if err = result.HasValidCheckDigits(); err != nil {
		return result, err
	}

	result.WarnOnUnknownBankCode()

	return result, nil
}

// CzechIbanFromNationalFormat parse Czech national format (e.g. 19-2000145399/0800) into IBAN structure for further
// processing.
func CzechIbanFromNationalFormat(text string) (CzechIBAN, error) {
	result := CzechIBAN{}

	accountNumber, bankCode, found := strings.Cut(text, "/")
	if !found {
		return result, fmt.Errorf("%w: %q", ErrCzNoSepBetweenAccountNumberAndBankCode, text)
	}

	accountPrefix, accountBase, found := strings.Cut(accountNumber, "-")
	if !found {
		accountBase, accountPrefix = accountPrefix, "0"
	}

	prefix, base, code, err := parseEachPart(accountPrefix, accountBase, bankCode)
	if err != nil {
		return result, err
	}

	result.accountNumberPrefix = prefix
	result.accountNumberBase = base
	result.bankCode = code
	result.ResetCheckDigits()

	remainder, err := CalculateCheckDigitsRemainder(result.MachineFormString())
	if err != nil {
		return result, err
	}

	result.checkDigits = 98 - remainder

	if err = result.HasValidCheckDigits(); err != nil { // duplicate, but better for double-check
		return result, err
	}

	result.WarnOnUnknownBankCode()

	return result, nil
}

func (i CzechIBAN) ResetCheckDigits() {
	i.checkDigits = 0
}

func (i CzechIBAN) HasValidCheckDigits() error {
	valid, err := ValidateCheckDigits(i.MachineFormString())
	switch {
	case err != nil:
		return err

	case !valid:
		return ErrWrongCheckDigits

	default:
		return nil
	}
}

func (i CzechIBAN) WarnOnUnknownBankCode() {
	if _, found := czechKnownBankCodes[i.bankCode]; !found {
		// TODO: log warning
	}
}

func (i CzechIBAN) CountryCode() string {
	return "CZ"
}

func (i CzechIBAN) CheckDigits() string {
	return fmt.Sprintf("%02d", i.checkDigits)
}

func (i CzechIBAN) BankCode() string {
	return fmt.Sprintf("%04d", i.bankCode)
}

func (i CzechIBAN) BankName() string {
	return czechKnownBankCodes[i.bankCode]
}

func (i CzechIBAN) AccountNumberPrefix() string {
	return fmt.Sprintf("%06d", i.accountNumberPrefix)
}

func (i CzechIBAN) AccountNumberBase() string {
	return fmt.Sprintf("%010d", i.accountNumberBase)
}

func (i CzechIBAN) FullCzechNationalString() string {
	return i.AccountNumberPrefix() + "-" + i.AccountNumberBase() + "/" + i.BankCode()
}

func (i CzechIBAN) CzechNationalString() string {
	if i.accountNumberPrefix == 0 {
		return fmt.Sprintf("%d/%s", i.accountNumberBase, i.BankCode())
	}

	return fmt.Sprintf("%d-%d/%s", i.accountNumberPrefix, i.accountNumberBase, i.BankCode())
}

func (i CzechIBAN) MachineFormString() string {
	return fmt.Sprintf(
		"%2.2s%02d%04d%06d%010d",
		i.CountryCode(),
		i.checkDigits,
		i.bankCode,
		i.accountNumberPrefix,
		i.accountNumberBase,
	)
}

func (i CzechIBAN) HumanFormString() string {
	return stringutils.SeparateGroups(i.MachineFormString(), 4, " ")
}

func checkAccountNumberPrefix(n uint64) error {
	switch {
	case n > maxValueOfAccountPrefix:
		return fmt.Errorf("%w: '%d'", ErrCzTooLargeAccountNumberPrefix, n)

	default:
		return nil
	}
}

func checkAccountNumberBase(n uint64) error {
	switch {
	case n > maxValueOfAccountBase:
		return fmt.Errorf("%w: '%d'", ErrCzTooLargeAccountNumberBase, n)

	case n == 0:
		return fmt.Errorf("%w: '%d'", ErrCzZeroAccountNumberBase, n)

	default:
		return nil
	}
}

func checkBankCode(n uint64) error {
	switch {
	case n > maxValueOfBankCode:
		return fmt.Errorf("%w: '%d'", ErrCzTooLargeBankCode, n)

	case n == 0:
		return fmt.Errorf("%w: '%d'", ErrCzZeroBankCode, n)

	default:
		return nil
	}
}

func parseEachPart(
	accountNumberPrefix string,
	accountNumberBase string,
	bankCode string,
) (prefix uint32, base uint32, code uint16, err error) {
	n, err := strconv.ParseUint(strings.TrimSpace(accountNumberPrefix), 10, 64)
	if err != nil {
		return prefix, base, code, fmt.Errorf("%w: %q", ErrCzParseAccountNumberPrefix, accountNumberPrefix)
	}
	if err = checkAccountNumberPrefix(n); err != nil {
		return prefix, base, code, err
	}
	prefix = uint32(n)

	n, err = strconv.ParseUint(strings.TrimSpace(accountNumberBase), 10, 64)
	if err != nil {
		return prefix, base, code, fmt.Errorf("%w: %q", ErrCzParseAccountNumberBase, accountNumberBase)
	}
	if err = checkAccountNumberBase(n); err != nil {
		return prefix, base, code, err
	}
	base = uint32(n)

	n, err = strconv.ParseUint(strings.TrimSpace(bankCode), 10, 64)
	if err != nil {
		return prefix, base, code, fmt.Errorf("%w: %q", ErrCzParseBankCode, bankCode)
	}
	if err = checkBankCode(n); err != nil {
		return prefix, base, code, err
	}
	code = uint16(n)

	return prefix, base, code, nil
}
