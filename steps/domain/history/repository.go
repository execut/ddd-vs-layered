package history

import (
    "context"

    "effective-architecture/steps/domain"
)

type IRepository interface {
    Create(ctx context.Context, history History) error
    List(ctx context.Context, id domain.LabelTemplateID) ([]History, error)
    Count(ctx context.Context, id domain.LabelTemplateID) (int, error)
    Cleanup(ctx context.Context, id domain.LabelTemplateID) error
}
