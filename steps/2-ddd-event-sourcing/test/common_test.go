package test_test

import "effective-architecture/steps/2-ddd-event-sourcing/domain"

const (
    expectedUUIDValue                         = "123e4567-e89b-12d3-a456-426655440000"
    expectedManufacturerOrganizationNameValue = "test manufacturer organization name"
)

var (
    expectedUUID, _                         = domain.NewLabelTemplateID(expectedUUIDValue)
    expectedManufacturerOrganizationName, _ = domain.NewManufacturerOrganizationName(
        expectedManufacturerOrganizationNameValue)
)
