package test

import (
    "context"
    "testing"

    "effective-architecture/steps/1-layered/labels"
    "github.com/stretchr/testify/require"
)

func TestLabels_Live(t *testing.T) {
    t.Parallel()
    var service *labels.Service
    repository, err := labels.NewRepository(t.Context())
    require.NoError(t, err)
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
}
