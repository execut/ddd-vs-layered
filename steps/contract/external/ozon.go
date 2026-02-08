//go:generate mockgen -package external -destination=ozon_mocks.go . IExternalServiceOzon
package external

import (
    "context"
    "errors"
)

var (
    ErrSkuNotFound = errors.New("sku not found")
)

type IExternalServiceOzon interface {
    Product(ctx context.Context, sku int64) (Product, error)
}

type Product struct {
    Category CategoryWithType
}

type CategoryWithType struct {
    Category Category
    TypeID   int64
}

type Category struct {
    ID             int64
    ParentCategory *Category
}
