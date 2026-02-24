package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"effective-architecture/steps/domain"
	"effective-architecture/steps/infrastructure/index/category"

	"github.com/jackc/pgx/v5"
)

var _ domain.IRepository = (*Repository)(nil)

type Repository struct {
	db                           *EventsRepository
	categoryVsTemplateRepository *category.Repository
	dispatcher                   *domain.Dispatcher
}

func NewRepository(db *EventsRepository,
	categoryVsTemplateRepository *category.Repository,
	dispatcher *domain.Dispatcher) *Repository {
	return &Repository{
		db:                           db,
		categoryVsTemplateRepository: categoryVsTemplateRepository,
		dispatcher:                   dispatcher,
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
		case "domain.LabelTemplateCategoryListUnlinkedEvent":
			err = applyEvent[domain.LabelTemplateCategoryListUnlinkedEvent](model, aggregate)
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

	err = r.dispatcher.Dispatch(ctx, aggregate)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) LoadByCategoryList(ctx context.Context,
	categoryList []domain.Category) (*domain.LabelTemplate, error) {
	for _, cat := range categoryList {
		aggregateID, err := r.categoryVsTemplateRepository.AggregateIDByCategory(ctx, cat)
		if err != nil {
			if errors.Is(err, category.ErrNotFound) {
				continue
			}

			return nil, err
		}

		if aggregateID == nil {
			continue
		}

		domainAggregateID, err := domain.NewLabelTemplateID(*aggregateID)
		if err != nil {
			return nil, err
		}

		aggregate, err := domain.NewLabelTemplate(domainAggregateID)
		if err != nil {
			return nil, err
		}

		err = r.Load(ctx, aggregate)
		if err != nil {
			return nil, err
		}

		return aggregate, nil
	}

	return nil, domain.ErrLabelTemplateForCategoryNotFound
}

func (r Repository) Cleanup(ctx context.Context, labelTemplateID domain.LabelTemplateID) error {
	err := r.categoryVsTemplateRepository.Cleanup(ctx, labelTemplateID)
	if err != nil {
		return err
	}

	return r.db.Cleanup(ctx, labelTemplateID.UUID)
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
