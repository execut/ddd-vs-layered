package domain

import (
	"context"
	"errors"
)

var (
	ErrLabelTemplateForCategoryNotFound = errors.New("шаблон этикетки для SKU не найден")
	ErrLabelAlreadyExists               = errors.New("генерация этикетки с таким идентификатором уже существует")
	ErrUnsupportedEventType             = errors.New("unsupported event type")
)

type Label struct {
	ID         LabelID
	Status     LabelStatus
	TemplateID LabelTemplateID
	SKU        int64
	Product    Product

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

	categoryList, _, err := l.ozonService.ProductData(ctx, sku)
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

func (l *Label) FillData(ctx context.Context) error {
	categoryList, product, err := l.ozonService.ProductData(ctx, l.SKU)
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

	err = l.addAndApplyEvent(LabelDataFilledEvent{
		LabelTemplateID: labelTemplate.ID,
		Product:         product,
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *Label) ApplyEvent(event LabelEvent) error {
	switch payload := event.(type) {
	case LabelGenerationStartedEvent:
		l.Status = LabelStatusGeneration
		l.SKU = payload.SKU
	case LabelDataFilledEvent:
		l.Status = LabelStatusDataFilled
		l.Product = payload.Product
		l.TemplateID = payload.LabelTemplateID
	default:
		return ErrUnsupportedEventType
	}

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
