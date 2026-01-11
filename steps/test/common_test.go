package test_test

import (
    "effective-architecture/steps/domain"
)

const (
    expectedUUIDValue                            = "123e4567-e89b-12d3-a456-426655440000"
    expectedManufacturerOrganizationNameValue    = "test manufacturer organization name"
    expectedNewManufacturerOrganizationNameValue = "new test manufacturer organization name"
    expectedManufacturerOrganizationAddressValue = "test manufacturer organization address"
    expectedManufacturerEmailValue               = "test@test.com"
    expectedManufacturerSiteValue                = "test.com"
)

var (
    expectedUUID, _                         = domain.NewLabelTemplateID(expectedUUIDValue)
    expectedManufacturerOrganizationName, _ = domain.NewOrganizationName(
        expectedManufacturerOrganizationNameValue)
    expectedManufacturerOrganizationAddress, _ = domain.NewOrganizationAddress(
        expectedManufacturerOrganizationAddressValue,
    )
    expectedManufacturerEmail, _ = domain.NewEmail(expectedManufacturerEmailValue)
    expectedManufacturerSite, _  = domain.NewSite(expectedManufacturerSiteValue)
    expectedManufacturer, _      = domain.NewManufacturer(
        expectedManufacturerOrganizationName,
        &expectedManufacturerOrganizationAddress,
        &expectedManufacturerSite,
        &expectedManufacturerEmail,
    )
    expectedNewManufacturerOrganizationName, _ = domain.NewOrganizationName(
        expectedNewManufacturerOrganizationNameValue)
    expectedNewManufacturer, _ = domain.NewManufacturer(
        expectedNewManufacturerOrganizationName,
        nil,
        nil,
        nil,
    )
)
