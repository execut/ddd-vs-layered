package domain

type OrganizationAddress string

func NewOrganizationAddress(address string) (OrganizationAddress, error) {
    return OrganizationAddress(address), nil
}
