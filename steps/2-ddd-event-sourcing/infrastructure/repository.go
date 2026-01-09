package infrastructure

import (
    "context"
    "errors"

    "effective-architecture/steps/2-ddd-event-sourcing/domain"
    "github.com/jackc/pgx/v5"
)

type Repository struct {
    db *LabelTemplateRepository
}

func NewRepository(db *LabelTemplateRepository) *Repository {
    return &Repository{
        db: db,
    }
}

func (r Repository) Load(ctx context.Context, aggregate *domain.LabelTemplate) error {
    model, err := r.db.Load(ctx, aggregate.ID.UUID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil
        }

        return err
    }

    domainManufacturerOrganizationName, err := domain.NewManufacturerOrganizationName(
        model.ManufacturerOrganizationName)
    if err != nil {
        return err
    }

    aggregate.ManufacturerOrganizationName = domainManufacturerOrganizationName

    return nil
}

func (r Repository) Save(ctx context.Context, aggregate *domain.LabelTemplate) error {
    model := LabelTemplate{
        ID:                           aggregate.ID.UUID.String(),
        ManufacturerOrganizationName: aggregate.ManufacturerOrganizationName.Name,
    }

    err := r.db.Create(ctx, model)
    if err != nil {
        return err
    }

    return nil
}
