package test_test

import (
    "context"
    "testing"

    labels2 "effective-architecture/steps/labels"
    "github.com/stretchr/testify/require"
)

func TestLabels_Live(t *testing.T) {
    t.Parallel()

    var service *labels2.Service

    repository, err := labels2.NewRepository(t.Context())
    require.NoError(t, err)
    t.Cleanup(func() {
        _ = repository.Truncate(context.Background())
    })

    t.Run("New", func(t *testing.T) {
        service = labels2.NewService(repository)
        require.NotNil(t, service)
    })

    t.Run("CreateLabelTemplate", func(t *testing.T) {
        err := service.CreateLabelTemplate(t.Context(), testUUID, testManufacturerOrganizationName)

        require.NoError(t, err)
    })
}
