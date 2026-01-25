package domain

import "context"

type IRepository interface {
    Load(ctx context.Context, aggregate *LabelTemplate) error
    Save(ctx context.Context, aggregate *LabelTemplate) error
    Cleanup(ctx context.Context, id LabelTemplateID) error
}
