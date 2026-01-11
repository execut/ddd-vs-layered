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

type Email string

func NewEmail(value string) (Email, error) {
    if len(value) > 255 || len(value) == 0 {
        return "", ErrEmailWrongLen
    }

    _, err := mail.ParseAddress(value)
    if err != nil {
        return "", ErrEmailWrongFormat
    }

    return Email(value), nil
}
