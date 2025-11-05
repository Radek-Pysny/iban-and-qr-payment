package iban

import (
	"fmt"
	"iter"
	"unicode"
)

func ValidateCheckDigits(iban string) (bool, error) {
	remainder, err := CalculateCheckDigitsRemainder(iban)

	return remainder == 1, err
}

func CalculateCheckDigitsRemainder(iban string) (uint8, error) {
	const nineDigitLimit = 100_000_000
	next, stop := iter.Pull2(IterateIbanForCheckDigits(iban))
	defer stop()

	var acc uint32
	for {
		digit, digitOK, exists := next()
		if !exists {
			break
		}
		if !digitOK {
			return 0, fmt.Errorf("%w: %d", ErrWrongCharacter, digit)
		}

		acc *= 10
		acc += uint32(digit)
		if acc >= nineDigitLimit {
			acc = acc % 97
		}
	}

	return uint8(acc % 97), nil
}

// IterateIbanCheckDigitsReminderStepByStep is not really meant to be used in production, it is included due to
// teaching and debugging purpose of piece-wise manner modulo-97 operation on a bigger unsigned integers
// (e.g. 219-bit).
func IterateIbanCheckDigitsReminderStepByStep(iban string) iter.Seq2[uint32, error] {
	return func(yield func(uint32, error) bool) {
		const nineDigitLimit = 100_000_000
		next, stop := iter.Pull2(IterateIbanForCheckDigits(iban))
		defer stop()

		var acc uint32
		for {
			digit, digitOK, exists := next()
			if !exists {
				break
			}
			if !digitOK {
				yield(0, fmt.Errorf("%w: %d", ErrWrongCharacter, digit))
				return
			}

			acc *= 10
			acc += uint32(digit)
			if acc >= nineDigitLimit {
				if !yield(acc, nil) {
					return
				}
				acc = acc % 97
			}
		}

		if !yield(acc, nil) {
			return
		}

		if acc < 97 {
			return
		}

		yield(acc%97, nil)
	}
}

// IterateIbanForCheckDigits is custom iterator that makes calculation of IBAN check digits simpler in readability.
func IterateIbanForCheckDigits(iban string) iter.Seq2[rune, bool] {
	return func(yield func(rune, bool) bool) {
		// Detection of where we should start
		count := 0
		index := 0
		for _, char := range iban {
			index++

			switch {
			case unicode.IsSpace(char):
				// skip whitespace
				break

			case char >= 'A' && char <= 'Z',
				char >= '0' && char <= '9':
				// count one more valid character
				count++

			default:
				// signalize error (returning codepoint of wrong character) and stop processing
				yield(char, false)
				return
			}

			if count == 4 { // got passed first 4 valid characters (national code and check digits)
				break // we already know index, from where to start processing
			}
		}

		process := func(input string) bool {
			for _, char := range input {
				switch {
				case unicode.IsSpace(char):
					continue

				case char >= '0' && char <= '9': // decode digit as just one digit
					if !yield(char-'0', true) {
						return false
					}

				case char >= 'A' && char <= 'Z': // decode letter as two separate digits
					x := char - 'A' + 10
					if !yield(x/10, true) {
						return false
					}
					if !yield(x%10, true) {
						return false
					}

				default:
					// signalize error (returning codepoint of wrong character) and stop processing
					yield(char, false)
					return false
				}
			}

			return true
		}

		// start processing after the "IBAN header"
		if !process(iban[index:]) {
			return
		}

		// processing of "IBAN header" at the end
		process(iban[:index])
	}
}
