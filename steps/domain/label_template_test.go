package domain_test

import (
    "testing"

    "effective-architecture/steps/domain"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLabelTemplate_Live(t *testing.T) {
    t.Parallel()

    var (
        labelTemplate *domain.LabelTemplate
        err           error
    )

    t.Run("New", func(t *testing.T) {
        labelTemplate, err = domain.NewLabelTemplate(expectedUUID)

        require.NoError(t, err)
        assert.Equal(t, expectedUUID, labelTemplate.ID)
    })

    t.Run("Create", func(t *testing.T) {
        err := labelTemplate.Create(expectedManufacturer)

        require.NoError(t, err)
        assert.Equal(t, expectedManufacturer, labelTemplate.Manufacturer)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
        err := labelTemplate.Create(expectedManufacturer)

        require.ErrorIs(t, err, domain.ErrLabelTemplateAlreadyCreated)
    })

    t.Run("Обновлять данные шаблона", func(t *testing.T) {
        err := labelTemplate.Update(expectedNewManufacturer)

        require.NoError(t, err)
        assert.Equal(t, expectedNewManufacturer, labelTemplate.Manufacturer)
    })

    t.Run("Удалять шаблон этикетки товара по UUID", func(t *testing.T) {
        err := labelTemplate.Delete()

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
        err := labelTemplate.Delete()

        require.ErrorIs(t, err, domain.ErrLabelTemplateAlreadyDeleted)
    })

    t.Run("CreateAgain", func(t *testing.T) {
        err := labelTemplate.Create(expectedManufacturer)

        require.NoError(t, err)
        assert.Equal(t, expectedManufacturer, labelTemplate.Manufacturer)
    })
}
