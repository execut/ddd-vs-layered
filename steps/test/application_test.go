package test_test

import (
    "context"
    "testing"

    "effective-architecture/steps/application"
    "effective-architecture/steps/infrastructure"
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
        app, err = application.NewApplication(repository)

        require.NoError(t, err)
        require.NotNil(t, app)
    })

    t.Run("CreateLabelTemplate", func(t *testing.T) {
        err := app.CreateLabelTemplate(t.Context(), expectedUUIDValue, expectedManufacturerOrganizationNameValue)

        require.NoError(t, err)
    })

    t.Run("GetLabelTemplate", func(t *testing.T) {
        result, err := app.GetLabelTemplate(t.Context(), expectedUUIDValue)

        require.NoError(t, err)
        assert.JSONEq(t, `{"id":"123e4567-e89b-12d3-a456-426655440000",`+
            `"manufacturerOrganizationName":"test manufacturer organization name"}`, result)
    })

    t.Run("DeleteLabelTemplate", func(t *testing.T) {
        err := app.DeleteLabelTemplate(t.Context(), expectedUUIDValue)

        require.NoError(t, err)
    })
}
