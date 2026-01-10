package test_test

import (
    "context"
    "strings"
    "testing"

    "effective-architecture/steps/labels"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestLabels_Live(t *testing.T) {
    t.Parallel()

    var service *labels.Service

    repository, err := labels.NewRepository(t.Context())
    require.NoError(t, err)

    _ = repository.Truncate(context.Background())

    t.Cleanup(func() {
        _ = repository.Truncate(context.Background())
    })

    t.Run("New", func(t *testing.T) {
        service = labels.NewService(repository)
        require.NotNil(t, service)
    })

    t.Run("CreateLabelTemplate", func(t *testing.T) {
        err := service.CreateLabelTemplate(t.Context(), testUUID, testManufacturerOrganizationName)

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
        err := service.CreateLabelTemplate(t.Context(), testUUID, testManufacturerOrganizationName)

        require.ErrorIs(t, err, labels.ErrLabelTemplateAlreadyCreated)
    })

    t.Run("GetLabelTemplate", func(t *testing.T) {
        result, err := service.GetLabelTemplate(t.Context(), testUUID)

        require.NoError(t, err)
        assert.JSONEq(t, `{"id":"123e4567-e89b-12d3-a456-426655440000","manufacturerOrganizationName":`+
            `"test manufacturer organization name"}`, result)
    })

    t.Run("DeleteLabelTemplate", func(t *testing.T) {
        err := service.DeleteLabelTemplate(t.Context(), testUUID)

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
        err := service.DeleteLabelTemplate(t.Context(), testUUID)

        require.ErrorIs(t, err, labels.ErrLabelTemplateAlreadyDeleted)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать шаблон, "+
        "если длина Наименования организации производителя", func(t *testing.T) {
        t.Run("> 255", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), testUUID, strings.Repeat("a", 256))

            require.ErrorIs(t, err, labels.ErrLabelTemplateWrongManufacturerOrganizationName)
        })
        t.Run("= 0", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), testUUID, "")

            require.ErrorIs(t, err, labels.ErrLabelTemplateWrongManufacturerOrganizationName)
        })
    })
}
