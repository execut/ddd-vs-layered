package domain

import (
    "errors"
    "net/url"
)

var (
    ErrSiteWrongLen = errors.New(
        "сайт должен быть до 255 символов в длину")
    ErrSiteWrongFormat = errors.New("сайт имеет не корректный формат")
)

type Site struct {
    Value string
}

func NewSite(value string) (Site, error) {
    if len(value) > 255 || len(value) == 0 {
        return Site{}, ErrSiteWrongLen
    }

    _, err := url.ParseRequestURI(value)
    if err != nil {
        return Site{}, ErrSiteWrongFormat
    }

    return Site{
        Value: value,
    }, nil
}
