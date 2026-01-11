package domain

type Manufacturer struct {
    OrganizationName    OrganizationName
    OrganizationAddress *OrganizationAddress
    Site                *Site
    Email               *Email
}

func NewManufacturer(organizationName OrganizationName, organizationAddress *OrganizationAddress,
    manufacturerSite *Site, manufacturerEmail *Email) (Manufacturer, error) {
    return Manufacturer{
        OrganizationName:    organizationName,
        OrganizationAddress: organizationAddress,
        Site:                manufacturerSite,
        Email:               manufacturerEmail,
    }, nil
}
