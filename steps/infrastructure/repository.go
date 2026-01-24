package infrastructure

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "time"

    "effective-architecture/steps/domain"
    "github.com/jackc/pgx/v5"
)

type Repository struct {
    db *EventsRepository
}

func NewRepository(db *EventsRepository) *Repository {
    return &Repository{
        db: db,
    }
}

func (r Repository) Load(ctx context.Context, aggregate *domain.LabelTemplate) error { //nolint:cyclop
    modelList, err := r.db.Load(ctx, aggregate.ID.UUID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil
        }

        return err
    }

    for _, model := range modelList {
        switch model.Type {
        case "domain.LabelTemplateCreatedEvent":
            err = applyEvent[domain.LabelTemplateCreatedEvent](model, aggregate)
            if err != nil {
                return err
            }
        case "domain.LabelTemplateDeletedEvent":
            err = applyEvent[domain.LabelTemplateDeletedEvent](model, aggregate)
            if err != nil {
                return err
            }
        case "domain.LabelTemplateUpdatedEvent":
            err = applyEvent[domain.LabelTemplateUpdatedEvent](model, aggregate)
            if err != nil {
                return err
            }
        case "domain.LabelTemplateCategoryListAddedEvent":
            err = applyEvent[domain.LabelTemplateCategoryListAddedEvent](model, aggregate)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func (r Repository) Save(ctx context.Context, aggregate *domain.LabelTemplate) error {
    eventModelList := make([]EventModel, 0, len(aggregate.Events))

    for _, event := range aggregate.Events {
        payload, err := json.Marshal(event)
        if err != nil {
            return err
        }

        eventModel := EventModel{
            Type:        fmt.Sprintf("%T", event),
            AggregateID: aggregate.ID.UUID.String(),
            Payload:     payload,
            CreatedAt:   time.Now(),
        }

        eventModelList = append(eventModelList, eventModel)
    }

    err := r.db.Save(ctx, eventModelList)
    if err != nil {
        return err
    }

    return nil
}

func applyEvent[T domain.LabelTemplateEvent](model EventModel, aggregate *domain.LabelTemplate) error {
    var event T

    err := json.Unmarshal(model.Payload, &event)
    if err != nil {
        return err
    }

    err = aggregate.ApplyEvent(event)
    if err != nil {
        return err
    }

    return nil
}
