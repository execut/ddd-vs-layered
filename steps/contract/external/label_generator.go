//go:generate mockgen -package external -destination=label_generator_mocks.go . ILabelGenerator
package external

import (
	"context"
	"effective-architecture/steps/contract"
)

type ILabelGenerator interface {
	Generate(ctx context.Context, product contract.Product) (LabelGeneratorFile, error)
}

type LabelGeneratorFile struct {
	Path string
}
