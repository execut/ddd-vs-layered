package test_test

import (
    "context"
    "testing"

    "effective-architecture/steps/2-ddd-event-sourcing/application"
    "effective-architecture/steps/2-ddd-event-sourcing/infrastructure"
    "github.com/stretchr/testify/require"
)

func TestApplication_Live(t *testing.T) {
    t.Parallel()

    var (
        app             *application.Application
        repository, err = infrastructure.NewLabelTemplateRepository()
    )
    if err != nil {
        panic(err)
    }

    t.Cleanup(func() {
        err = repository.Truncate(context.Background())
        _ = err
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
}
