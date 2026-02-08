//go:generate mockgen -package domain -destination=ozon_service_mock.go . IServiceOzon
package domain

import "context"

type IServiceOzon interface {
    CategoryList(ctx context.Context, sku int64) ([]Category, error)
}
