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
	ID           LabelID
	Status       LabelStatus
	TemplateID   LabelTemplateID
	SKU          int64
	Product      Product
	Manufacturer Manufacturer
	File         *LabelFile

	ozonService     IServiceOzon
	labelGenerator  ILabelGenerator
	labelRepository IRepository
	Events          []LabelEvent
}

func NewLabel(id LabelID, ozonService IServiceOzon, labelRepository IRepository,
	labelGenerator ILabelGenerator) *Label {
	return &Label{
		ID:              id,
		ozonService:     ozonService,
		labelRepository: labelRepository,
		labelGenerator:  labelGenerator,
	}
}

func (l *Label) StartGeneration(sku int64) error {
	if l.Status == LabelStatusGeneration {
		return ErrLabelAlreadyExists
	}

	err := l.addAndApplyEvent(LabelGenerationStartedEvent{
		SKU: sku,
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
		Manufacturer:    labelTemplate.Manufacturer,
	})
	if err != nil {
		return err
	}

	return nil
}

func (l *Label) Generate(ctx context.Context) error {
	labelFile, err := l.labelGenerator.Generate(ctx, l.Product, l.Manufacturer, l.SKU)
	if err != nil {
		return err
	}

	err = l.addAndApplyEvent(LabelGeneratedEvent{
		LabelFile: labelFile,
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
		l.Manufacturer = payload.Manufacturer
	case LabelGeneratedEvent:
		l.Status = LabelStatusGenerated
		l.File = &payload.LabelFile
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
