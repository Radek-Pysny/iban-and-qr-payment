package qrpay

type optionalKey func(*ShortPaymentDescriptor) error

func New(iban string, bic string, keys ...optionalKey) (*ShortPaymentDescriptor, error) {
	account, err := buildAccount(iban, bic)
	if err != nil {
		return nil, err
	}

	result := ShortPaymentDescriptor{
		Account: account,
	}

	for _, key := range keys {
		if err = key(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func buildAccount(iban string, bic string) (string, error) {
	// TODO: check iban
	result := iban

	if len(bic) > 0 {
		result += "+" + bic
	}

	return result, nil
}
