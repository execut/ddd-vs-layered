//go:generate mockgen -package domain -destination=repository_mock.go . IRepository
package domain

import (
    "context"
)

type IRepository interface {
    Load(ctx context.Context, aggregate *LabelTemplate) error
    Save(ctx context.Context, aggregate *LabelTemplate) error
    Cleanup(ctx context.Context, id LabelTemplateID) error
    LoadByCategoryList(ctx context.Context, categoryList []Category) (*LabelTemplate, error)
}
