package domain

import "context"

type ILabelRepository interface {
    Load(ctx context.Context, aggregate *Label) error
    Save(ctx context.Context, aggregate *Label) error
}
