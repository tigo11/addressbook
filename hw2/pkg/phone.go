package pkg

import (
	"errors"
	"regexp"
	"strings"
)

var ErrInvalidPhoneNumber = errors.New("invalid phone number")

func PhoneNormalize(phone string) (normalizedPhone string, err error) {
	re := regexp.MustCompile("[^0-9]")
	normalizedPhone = re.ReplaceAllString(phone, "")

	// Если номер начинается с "8", заменяем на "+7"
	if len(normalizedPhone) == 11 {
		normalizedPhone = "+7" + normalizedPhone[len(normalizedPhone)-10:]
	}

	// Если номер не начинается ни с "+7", ни с "8", считаем его некорректным
	if !strings.HasPrefix(normalizedPhone, "+7") {
		return "", ErrInvalidPhoneNumber
	}

	if len(normalizedPhone) > 12 || len(normalizedPhone) < 11 {
		return "", ErrInvalidPhoneNumber
	}
	return normalizedPhone, nil
}