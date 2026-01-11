package test_test

import "effective-architecture/steps/labels"

const (
    expectedUUID                               = "123e4567-e89b-12d3-a456-426655440000"
    expectedManufacturerOrganizationName       = "test manufacturer organization name"
    expectedManufacturerOrganizationAddress    = "test manufacturer organization address"
    expectedManufacturerSite                   = "test.com"
    expectedManufacturerEmail                  = "email@test.com"
    expectedNewManufacturerOrganizationName    = "new test manufacturer organization name"
    expectedNewManufacturerOrganizationAddress = "new test manufacturer organization address"
    expectedNewManufacturerSite                = "new-test.com"
    expectedNewManufacturerEmail               = "new-email@test.com"
)

var (
    expectedManufacturer = labels.Manufacturer{
        OrganizationName: expectedManufacturerOrganizationName,
    }
    expectedManufacturerWithAllFields = labels.Manufacturer{
        OrganizationName:    expectedManufacturerOrganizationName,
        OrganizationAddress: expectedManufacturerOrganizationAddress,
        Email:               expectedManufacturerEmail,
        Site:                expectedManufacturerSite,
    }
    expectedNewManufacturer = labels.Manufacturer{
        OrganizationName: expectedNewManufacturerOrganizationName,
    }
    expectedNewManufacturerWithAllFields = labels.Manufacturer{
        OrganizationName:    expectedNewManufacturerOrganizationName,
        OrganizationAddress: expectedNewManufacturerOrganizationAddress,
        Email:               expectedNewManufacturerEmail,
        Site:                expectedNewManufacturerSite,
    }
)
