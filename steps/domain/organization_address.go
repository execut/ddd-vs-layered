package domain

import "errors"

var (
    ErrManufacturerOrganizationAddressWrongLen = errors.New(
        "адрес должен быть до 255 символов в длину")
)

type OrganizationAddress string

func NewOrganizationAddress(value string) (OrganizationAddress, error) {
    if len(value) > 255 || len(value) == 0 {
        return "", ErrManufacturerOrganizationAddressWrongLen
    }

    return OrganizationAddress(value), nil
}
