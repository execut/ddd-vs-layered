package e2e_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLabelLive(t *testing.T) {
    t.Parallel()
    t.Run("Создавать шаблон этикетки товара с UUID и Наименованием организации производителя", func(t *testing.T) {
        output, err := runBinary([]string{
            "labels-create-template",
            "--id", testUUID,
            "--manufacturer-organization-name", testManufacturerOrganizationName,
        })

        require.NoError(t, err, "output:"+output)
        assert.Equal(t, "1", output)
    })
    t.Run("Получать данные шаблона в JSON", func(t *testing.T) {
        output, err := runBinary([]string{
            "labels-get-template",
            "--id", testUUID,
        })

        require.NoError(t, err)
        assert.Equal(t, `{"id":"123e4567-e89b-12d3-a456-426655440000","manufacturerOrganizationName":"test manufacturer organization name"}`, output)
    })
    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
        out, err := runBinary([]string{
            "labels-create-template",
            "--id", testUUID,
            "--manufacturer-organization-name", testManufacturerOrganizationName,
        })

        require.Error(t, err)
        assert.Contains(t, out, "попытка создать уже существующий шаблон")
    })
    t.Run("Удалять шаблон этикетки товара по UUID", func(t *testing.T) {
        output, err := runBinary([]string{
            "labels-delete-template",
            "--id", testUUID,
        })

        require.NoError(t, err)
        assert.Equal(t, `1`, output)
    })
}
