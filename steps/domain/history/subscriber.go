package history

import (
    "context"

    "effective-architecture/steps/domain"
)

type Subscriber struct {
    repository IRepository
}

func NewSubscriber(repository IRepository) *Subscriber {
    return &Subscriber{repository: repository}
}

func (s Subscriber) Handle(ctx context.Context,
    aggregate *domain.LabelTemplate, event domain.LabelTemplateEvent) error {
    currentCount, err := s.repository.Count(ctx, aggregate.ID)
    if err != nil {
        return err
    }

    history, err := NewHistoryFromEvent(aggregate, event, currentCount)
    if err != nil {
        return err
    }

    err = s.repository.Create(ctx, history)
    if err != nil {
        return err
    }

    return nil
}
