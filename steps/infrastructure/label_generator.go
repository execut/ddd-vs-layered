package infrastructure

import (
	"context"
	"effective-architecture/steps/contract"

	"effective-architecture/steps/contract/external"
	"effective-architecture/steps/domain"
)

var _ domain.ILabelGenerator = (*LabelGenerator)(nil)

type LabelGenerator struct {
	externalLabelGenerator external.ILabelGenerator
}

func NewLabelGenerator(externalLabelGenerator external.ILabelGenerator) *LabelGenerator {
	return &LabelGenerator{externalLabelGenerator: externalLabelGenerator}
}

func (s *LabelGenerator) Generate(ctx context.Context, product domain.Product, manufacturer domain.Manufacturer,
	sku int64) (domain.LabelFile, error) {
	contractManufacturer := contract.Manufacturer{
		OrganizationName: manufacturer.OrganizationName.Name,
		Email:            manufacturer.Email.Value,
	}
	if manufacturer.Site != nil {
		contractManufacturer.Site = manufacturer.Site.Value
	}

	if manufacturer.OrganizationAddress != nil {
		contractManufacturer.OrganizationAddress = manufacturer.OrganizationAddress.Address
	}

	externalProduct := contract.Product{
		Name:         product.Name,
		SKU:          sku,
		Manufacturer: contractManufacturer,
	}

	result, err := s.externalLabelGenerator.Generate(ctx, externalProduct)
	if err != nil {
		return domain.LabelFile{}, err
	}

	labelFile, err := domain.NewLabelFile(result.Path)
	if err != nil {
		return domain.LabelFile{}, err
	}

	return labelFile, nil
}
