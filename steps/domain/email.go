package domain

import (
    "errors"
    "net/mail"
)

var (
    ErrEmailWrongLen = errors.New(
        "email должен быть до 255 символов в длину")
    ErrEmailWrongFormat = errors.New("email имеет не корректный формат")
)

type Email struct {
    Value string
}

func NewEmail(value string) (Email, error) {
    if len(value) > 255 || len(value) == 0 {
        return Email{}, ErrEmailWrongLen
    }

    _, err := mail.ParseAddress(value)
    if err != nil {
        return Email{}, ErrEmailWrongFormat
    }

    return Email{
        Value: value,
    }, nil
}
