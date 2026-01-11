package domain

type Email string

func NewEmail(address string) (Email, error) {
    return Email(address), nil
}
