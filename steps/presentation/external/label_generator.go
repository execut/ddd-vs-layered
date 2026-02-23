package external

import (
	"context"
	"effective-architecture/steps/contract"
	"effective-architecture/steps/contract/external"
)

type LabelGenerator struct {
}

func NewLabelGenerator() *LabelGenerator {
	return &LabelGenerator{}
}

func (l *LabelGenerator) Generate(ctx context.Context, product contract.Product) (external.LabelGeneratorFile, error) {
	_ = ctx
	_ = product

	return external.LabelGeneratorFile{}, nil
}
