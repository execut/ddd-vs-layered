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

var _ domain.ILabelRepository = (*LabelRepository)(nil)

type LabelRepository struct {
    db *EventsRepository
}

func NewLabelRepository(db *EventsRepository) *LabelRepository {
    return &LabelRepository{
        db: db,
    }
}

func (r LabelRepository) Load(ctx context.Context, aggregate *domain.Label) error {
    modelList, err := r.db.Load(ctx, aggregate.ID.UUID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil
        }

        return err
    }

    for _, model := range modelList {
        if model.Type == "domain.LabelGenerationStartedEvent" {
            err = applyLabelEvent[domain.LabelGenerationStartedEvent](model, aggregate)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func (r LabelRepository) Save(ctx context.Context, aggregate *domain.Label) error {
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

func applyLabelEvent[T domain.LabelTemplateEvent](model EventModel, aggregate *domain.Label) error {
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
