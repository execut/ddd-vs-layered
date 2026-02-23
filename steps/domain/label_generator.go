//go:generate mockgen -package domain -destination=label_generator_mock.go . ILabelGenerator
package domain

import (
	"context"
)

type ILabelGenerator interface {
	Generate(ctx context.Context, product Product, manufacturer Manufacturer, sku int64) (LabelFile, error)
}
