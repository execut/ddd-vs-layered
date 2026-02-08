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

func (s *ServiceOzon) CategoryList(ctx context.Context, sku int64) ([]domain.Category, error) {
    product, err := s.externalServiceOzon.Product(ctx, sku)
    if err != nil {
        return nil, err
    }

    category, err := domain.NewCategory(product.Category.Category.ID, &product.Category.TypeID)
    if err != nil {
        return nil, err
    }

    categoryList := []domain.Category{category}

    ozonCategory := product.Category.Category
    for {
        if ozonCategory.ParentCategory == nil {
            return categoryList, nil
        }

        ozonCategory = *ozonCategory.ParentCategory

        category, err = domain.NewCategory(ozonCategory.ID, nil)
        if err != nil {
            return nil, err
        }

        categoryList = append(categoryList, category)
    }
}
