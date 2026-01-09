package domain

type ManufacturerOrganizationName struct {
    Name string
}

func NewManufacturerOrganizationName(name string) (ManufacturerOrganizationName, error) {
    return ManufacturerOrganizationName{Name: name}, nil
}
