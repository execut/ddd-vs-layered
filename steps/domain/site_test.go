package test_test

import (
    "strings"
    "testing"

    "effective-architecture/steps/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewSite_errors(t *testing.T) {
    t.Parallel()
    t.Run("Ошибка, если длина", func(t *testing.T) {
        t.Parallel()
        t.Run(">255", func(t *testing.T) {
            _, err := domain.NewSite(strings.Repeat("a", 256))

            require.ErrorIs(t, err, domain.ErrSiteWrongLen)
        })
        t.Run("или =0", func(t *testing.T) {
            _, err := domain.NewSite("")

            require.ErrorIs(t, err, domain.ErrSiteWrongLen)
        })
    })

    t.Run("если формат не корректный", func(t *testing.T) {
        t.Parallel()

        _, err := domain.NewSite("test")

        assert.ErrorIs(t, err, domain.ErrSiteWrongFormat)
    })
}
