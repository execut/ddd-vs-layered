package infrastructure

import (
    "context"
    "errors"

    domain2 "effective-architecture/steps/domain"
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

func (r Repository) Load(ctx context.Context, aggregate *domain2.LabelTemplate) error {
    model, err := r.db.Load(ctx, aggregate.ID.UUID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil
        }

        return err
    }

    domainManufacturerOrganizationName, err := domain2.NewManufacturerOrganizationName(
        model.ManufacturerOrganizationName)
    if err != nil {
        return err
    }

    aggregate.ManufacturerOrganizationName = domainManufacturerOrganizationName

    return nil
}

func (r Repository) Save(ctx context.Context, aggregate *domain2.LabelTemplate) error {
    var (
        createEvent *domain2.LabelTemplateCreatedEvent
        deleteEvent *domain2.LabelTemplateDeletedEvent
    )

    for _, event := range aggregate.Events {
        switch payload := event.(type) {
        case domain2.LabelTemplateCreatedEvent:
            createEvent = &payload
        case domain2.LabelTemplateDeletedEvent:
            deleteEvent = &payload
        }
    }

    if createEvent != nil {
        model := LabelTemplate{
            ID:                           aggregate.ID.UUID.String(),
            ManufacturerOrganizationName: createEvent.ManufacturerOrganizationName.Name,
        }

        err := r.db.Create(ctx, model)
        if err != nil {
            return err
        }
    }

    if deleteEvent != nil {
        err := r.db.Delete(ctx, aggregate.ID.UUID.String())
        if err != nil {
            return err
        }
    }

    return nil
}
