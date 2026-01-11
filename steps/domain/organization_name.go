package domain

import "errors"

var (
    ErrManufacturerOrganizationNameWrongLen = errors.New(
        "название организации производителя должно быть до 255 символов в длину")
)

type OrganizationName struct {
    Name string
}

func NewOrganizationName(value string) (OrganizationName, error) {
    if len(value) > 255 || len(value) == 0 {
        return OrganizationName{}, ErrManufacturerOrganizationNameWrongLen
    }

    return OrganizationName{Name: value}, nil
}
