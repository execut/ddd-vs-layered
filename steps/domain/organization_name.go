package domain

import "errors"

var (
    ErrManufacturerOrganizationNameWrongLen = errors.New(
        "название организации производителя должно быть до 255 символов в длину")
)

type OrganizationName struct {
    Name string
}

func NewOrganizationName(name string) (OrganizationName, error) {
    if len(name) > 255 || len(name) == 0 {
        return OrganizationName{}, ErrManufacturerOrganizationNameWrongLen
    }

    return OrganizationName{Name: name}, nil
}
