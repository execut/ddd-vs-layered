package domain

type Site string

func NewSite(address string) (Site, error) {
    return Site(address), nil
}
