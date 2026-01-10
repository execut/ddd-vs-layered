package test_test

import (
    "testing"

    "effective-architecture/steps/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewLabelID(t *testing.T) {
    labelID, err := domain.NewLabelTemplateID(expectedUUIDValue)

    require.NoError(t, err)
    assert.Equal(t, expectedUUIDValue, labelID.UUID.String())
}
