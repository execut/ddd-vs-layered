package test_test

import (
    domain2 "effective-architecture/steps/domain"
)

const (
    expectedUUIDValue                         = "123e4567-e89b-12d3-a456-426655440000"
    expectedManufacturerOrganizationNameValue = "test manufacturer organization name"
)

var (
    expectedUUID, _                         = domain2.NewLabelTemplateID(expectedUUIDValue)
    expectedManufacturerOrganizationName, _ = domain2.NewManufacturerOrganizationName(
        expectedManufacturerOrganizationNameValue)
)
