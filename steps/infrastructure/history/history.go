package history

type History struct {
    ID                                      string
    OrderKey                                int
    Action                                  string
    NewManufacturerOrganizationNameValue    *string
    NewManufacturerOrganizationAddressValue *string
    NewManufacturerEmailValue               *string
    NewManufacturerSiteValue                *string
}

type Category struct {
    CategoryID int64
    TypeID     *int64
}
