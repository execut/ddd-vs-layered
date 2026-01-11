package test_test

import (
    "strings"
    "testing"

    "effective-architecture/steps/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewEmail_errors(t *testing.T) {
    t.Parallel()
    t.Run("Ошибка, если длина", func(t *testing.T) {
        t.Parallel()
        t.Run(">255", func(t *testing.T) {
            _, err := domain.NewEmail(strings.Repeat("a", 256-9) + "@test.com")

            require.ErrorIs(t, err, domain.ErrEmailWrongLen)
        })
        t.Run("или =0", func(t *testing.T) {
            _, err := domain.NewEmail("")

            require.ErrorIs(t, err, domain.ErrEmailWrongLen)
        })
    })

    t.Run("если формат не корректный", func(t *testing.T) {
        t.Parallel()

        _, err := domain.NewEmail("test")

        assert.ErrorIs(t, err, domain.ErrEmailWrongFormat)
    })
}
