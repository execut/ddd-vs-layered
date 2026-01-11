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
        err := service.CreateLabelTemplate(t.Context(), expectedUUID, expectedManufacturer)

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
        err := service.CreateLabelTemplate(t.Context(), expectedUUID, expectedManufacturer)

        require.ErrorIs(t, err, labels.ErrLabelTemplateAlreadyCreated)
    })

    t.Run("GetLabelTemplate", func(t *testing.T) {
        result, err := service.GetLabelTemplate(t.Context(), expectedUUID)

        require.NoError(t, err)
        assert.JSONEq(t, `{"id":"123e4567-e89b-12d3-a456-426655440000","manufacturer":{"organizationName":`+
            `"test manufacturer organization name"}}`, result)
    })

    t.Run("Обновлять данные шаблона", func(t *testing.T) {
        err = service.UpdateLabelTemplate(t.Context(), expectedUUID, expectedNewManufacturer)

        require.NoError(t, err)
        result, err := service.GetLabelTemplate(t.Context(), expectedUUID)

        require.NoError(t, err)
        assert.Contains(t, result, expectedNewManufacturerOrganizationName)

        t.Run("и не давать это делать при ошибках из предыдущих пунктов", func(t *testing.T) {
            err = service.UpdateLabelTemplate(t.Context(), expectedUUID, labels.Manufacturer{
                OrganizationName: strings.Repeat("a", 256),
            })

            require.ErrorIs(t, labels.ErrLabelTemplateWrongManufacturerOrganizationName, err)
        })
    })

    t.Run("DeleteLabelTemplate", func(t *testing.T) {
        err := service.DeleteLabelTemplate(t.Context(), expectedUUID)

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
        err := service.DeleteLabelTemplate(t.Context(), expectedUUID)

        require.ErrorIs(t, err, labels.ErrLabelTemplateAlreadyDeleted)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке создать шаблон, "+
        "если длина Наименования организации производителя", func(t *testing.T) {
        t.Run("> 255", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), expectedUUID, labels.Manufacturer{
                OrganizationName: strings.Repeat("a", 256),
            })

            require.ErrorIs(t, err, labels.ErrLabelTemplateWrongManufacturerOrganizationName)
        })
        t.Run("= 0", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), expectedUUID, labels.Manufacturer{
                OrganizationName: "",
            })

            require.ErrorIs(t, err, labels.ErrLabelTemplateWrongManufacturerOrganizationName)
        })
    })

    t.Run("Указывать и получать поля Адрес, Email, сайт", func(t *testing.T) {
        t.Run("при создании", func(t *testing.T) {
            err := service.CreateLabelTemplate(t.Context(), expectedUUID, expectedManufacturerWithAllFields)

            require.NoError(t, err)
            result, err := service.GetLabelTemplate(t.Context(), expectedUUID)
            require.NoError(t, err)
            assert.Contains(t, result, expectedManufacturerOrganizationAddress)
            assert.Contains(t, result, expectedManufacturerEmail)
            assert.Contains(t, result, expectedManufacturerSite)
        })
        t.Run("и обновлении", func(t *testing.T) {
            err := service.UpdateLabelTemplate(t.Context(), expectedUUID, expectedNewManufacturerWithAllFields)

            require.NoError(t, err)
            result, err := service.GetLabelTemplate(t.Context(), expectedUUID)
            require.NoError(t, err)
            assert.Contains(t, result, expectedNewManufacturerOrganizationAddress)
            assert.Contains(t, result, expectedNewManufacturerEmail)
            assert.Contains(t, result, expectedNewManufacturerSite)
        })
    })
}
