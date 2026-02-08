package external

import (
    "context"

    "effective-architecture/steps/contract/external"
)

type ServiceOzon struct{}

func NewServiceOzon() *ServiceOzon {
    return &ServiceOzon{}
}

func (ServiceOzon) Product(_ context.Context, _ int64) (external.Product, error) {
    return external.Product{}, nil
}
