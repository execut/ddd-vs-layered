package category

import (
    "context"

    "effective-architecture/steps/domain"
)

type Subscriber struct {
    repository *Repository
}

func NewSubscriber(repository *Repository) *Subscriber {
    return &Subscriber{repository: repository}
}

func (s Subscriber) Handle(ctx context.Context,
    aggregate *domain.LabelTemplate, event domain.LabelTemplateEvent) error {
    switch payload := event.(type) {
    case domain.LabelTemplateCategoryListAddedEvent:
        for _, category := range payload.CategoryList {
            index := Index{
                LabelTemplateID: aggregate.ID.UUID.String(),
                TypeID:          category.TypeID,
                CategoryID:      category.CategoryID,
            }

            err := s.repository.Create(ctx, index)
            if err != nil {
                return err
            }
        }
    case domain.LabelTemplateCategoryListUnlinkedEvent:
        for _, category := range payload.CategoryList {
            index := Index{
                LabelTemplateID: aggregate.ID.UUID.String(),
                TypeID:          category.TypeID,
                CategoryID:      category.CategoryID,
            }

            err := s.repository.Delete(ctx, index)
            if err != nil {
                return err
            }
        }
    case domain.LabelTemplateDeletedEvent:
        for _, category := range aggregate.CategoryList {
            index := Index{
                LabelTemplateID: aggregate.ID.UUID.String(),
                TypeID:          category.TypeID,
                CategoryID:      category.CategoryID,
            }

            err := s.repository.Delete(ctx, index)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
