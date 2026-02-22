package domain

import (
	"context"
	"errors"
)

var (
	ErrLabelTemplateForCategoryNotFound = errors.New("шаблон этикетки для SKU не найден")
	ErrLabelAlreadyExists               = errors.New("генерация этикетки с таким идентификатором уже существует")
)

type Label struct {
	ID     LabelID
	Status LabelStatus

	ozonService     IServiceOzon
	labelRepository IRepository
	Events          []LabelEvent
}

func NewLabel(id LabelID, ozonService IServiceOzon, labelRepository IRepository) *Label {
	return &Label{
		ID:              id,
		ozonService:     ozonService,
		labelRepository: labelRepository,
	}
}

func (l *Label) StartGeneration(ctx context.Context, sku int64) error {
	if l.Status == LabelStatusGeneration {
		return ErrLabelAlreadyExists
	}

	categoryList, err := l.ozonService.CategoryList(ctx, sku)
	if err != nil {
		return err
	}

	labelTemplate, err := l.labelRepository.LoadByCategoryList(ctx, categoryList)
	if err != nil {
		return err
	}

	if labelTemplate == nil {
		return ErrLabelTemplateForCategoryNotFound
	}

	err = l.addAndApplyEvent(LabelGenerationStartedEvent{
		LabelTemplateID: labelTemplate.ID,
		SKU:             sku,
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *Label) ApplyEvent(event LabelEvent) error {
	_, ok := event.(LabelGenerationStartedEvent)
	if !ok {
		return nil
	}

	l.Status = LabelStatusGeneration

	return nil
}

func (l *Label) addAndApplyEvent(event LabelEvent) error {
	l.Events = append(l.Events, event)

	err := l.ApplyEvent(event)
	if err != nil {
		return err
	}

	return nil
}
