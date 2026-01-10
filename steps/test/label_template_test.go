package test_test

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
        err := labelTemplate.Create(expectedManufacturerOrganizationName)

        require.NoError(t, err)
        assert.Equal(t, expectedManufacturerOrganizationName, labelTemplate.ManufacturerOrganizationName)
    })
}
