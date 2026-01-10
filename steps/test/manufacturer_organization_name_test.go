package test_test

import (
    "strings"
    "testing"

    "effective-architecture/steps/domain"
    "github.com/stretchr/testify/require"
)

func TestNewManufacturerOrganizationName_errors(t *testing.T) {
    t.Parallel()
    t.Run("Ошибка, если длина Наименования организации производителя", func(t *testing.T) {
        t.Parallel()
        t.Run(">255", func(t *testing.T) {
            _, err := domain.NewManufacturerOrganizationName(strings.Repeat("a", 256))

            require.ErrorIs(t, err, domain.ErrManufacturerOrganizationNameWrongLen)
        })
        t.Run("или =0", func(t *testing.T) {
            _, err := domain.NewManufacturerOrganizationName("")

            require.ErrorIs(t, err, domain.ErrManufacturerOrganizationNameWrongLen)
        })
    })
}
