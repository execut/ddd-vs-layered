package test_test

import (
    "context"
    "testing"

    "effective-architecture/steps/application"
    "effective-architecture/steps/domain"
    "effective-architecture/steps/infrastructure"
    "effective-architecture/steps/infrastructure/history"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestApplication_Live(t *testing.T) {
    t.Parallel()

    var (
        app             *application.Application
        repository, err = infrastructure.NewEventsRepository()
    )
    require.NoError(t, err)

    _ = repository.Truncate(context.Background())

    t.Cleanup(func() {
        _ = repository.Truncate(context.Background())
    })
    t.Run("New", func(t *testing.T) {
        historyRepository, err := history.NewRepository()
        if err != nil {
            panic(err)
        }

        app, err = application.NewApplication(repository, historyRepository)

        require.NoError(t, err)
        require.NotNil(t, app)
    })

    t.Run("CreateLabelTemplate", func(t *testing.T) {
        err := app.CreateLabelTemplate(t.Context(), expectedUUIDValue, application.Manufacturer{
            OrganizationName: expectedManufacturerOrganizationNameValue,
        })

        require.NoError(t, err)
    })

    t.Run("GetLabelTemplate", func(t *testing.T) {
        result, err := app.GetLabelTemplate(t.Context(), expectedUUIDValue)

        require.NoError(t, err)
        assert.JSONEq(t, `{"id":"123e4567-e89b-12d3-a456-426655440000",`+
            `"manufacturer":{"organizationName":"test manufacturer organization name"}}`, result)
    })

    t.Run("UpdateLabelTemplate", func(t *testing.T) {
        err := app.UpdateLabelTemplate(t.Context(), expectedUUIDValue, application.Manufacturer{
            OrganizationName:    expectedNewManufacturerOrganizationNameValue,
            OrganizationAddress: expectedManufacturerOrganizationAddressValue,
            Email:               expectedManufacturerEmailValue,
            Site:                expectedManufacturerSiteValue,
        })

        require.NoError(t, err)
        result, err := app.GetLabelTemplate(t.Context(), expectedUUIDValue)
        require.NoError(t, err)
        assert.Contains(t, result, expectedNewManufacturerOrganizationNameValue)
        assert.Contains(t, result, expectedManufacturerOrganizationAddressValue)
        assert.Contains(t, result, expectedManufacturerEmailValue)
        assert.Contains(t, result, expectedManufacturerSiteValue)
    })

    t.Run("Удалять шаблон этикетки товара по UUID", func(t *testing.T) {
        err := app.DeleteLabelTemplate(t.Context(), expectedUUIDValue)

        require.NoError(t, err)
    })

    t.Run("Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
        err := app.DeleteLabelTemplate(t.Context(), expectedUUIDValue)

        require.ErrorIs(t, err, domain.ErrLabelTemplateAlreadyDeleted)
    })

    t.Run("Чтобы писалась история операций над шаблонами с возможностью выводить все"+
        " данные в json", func(t *testing.T) {
        out, err := app.LabelTemplateHistoryList(t.Context(), expectedUUIDValue)

        require.NoError(t, err)
        assert.JSONEq(t, `
    [{
       "orderKey": 1,
       "action": "created",
       "newManufacturerOrganizationName": "test manufacturer organization name"
    },
    {
       "orderKey": 2,
       "action": "updated",
       "newManufacturerOrganizationName": "new test manufacturer organization name",
       "newManufacturerOrganizationAddress": "test manufacturer organization address",
       "newManufacturerEmail": "test@test.com",
       "newManufacturerSite": "https://test.com"
    },
    {
       "orderKey": 3,
       "action": "deleted"
    }]
    `, out)
    })
}
