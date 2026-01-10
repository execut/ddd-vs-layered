package domain

import "errors"

var (
    ErrManufacturerOrganizationNameWrongLen = errors.New(
        "название организации производителя должно быть до 255 символов в длину")
)

type ManufacturerOrganizationName struct {
    Name string
}

func NewManufacturerOrganizationName(name string) (ManufacturerOrganizationName, error) {
    if len(name) > 255 || len(name) == 0 {
        return ManufacturerOrganizationName{}, ErrManufacturerOrganizationNameWrongLen
    }

    return ManufacturerOrganizationName{Name: name}, nil
}
