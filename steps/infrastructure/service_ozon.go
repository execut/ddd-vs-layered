package infrastructure

import (
	"context"

	"effective-architecture/steps/contract/external"
	"effective-architecture/steps/domain"
)

var _ domain.IServiceOzon = (*ServiceOzon)(nil)

type ServiceOzon struct {
	externalServiceOzon external.IExternalServiceOzon
}

func NewServiceOzon(externalServiceOzon external.IExternalServiceOzon) *ServiceOzon {
	return &ServiceOzon{externalServiceOzon: externalServiceOzon}
}

func (s *ServiceOzon) ProductData(ctx context.Context, sku int64) ([]domain.Category, domain.Product, error) {
	product, err := s.externalServiceOzon.Product(ctx, sku)
	if err != nil {
		return nil, domain.Product{}, err
	}

	domainProduct, err := domain.NewProduct(product.Name)
	if err != nil {
		return nil, domain.Product{}, err
	}

	category, err := domain.NewCategory(product.Category.Category.ID, &product.Category.TypeID)
	if err != nil {
		return nil, domain.Product{}, err
	}

	categoryList := []domain.Category{category}

	ozonCategory := product.Category.Category
	for {
		if ozonCategory.ParentCategory == nil {
			return categoryList, domainProduct, nil
		}

		ozonCategory = *ozonCategory.ParentCategory

		category, err = domain.NewCategory(ozonCategory.ID, nil)
		if err != nil {
			return nil, domain.Product{}, err
		}

		categoryList = append(categoryList, category)
	}
}
