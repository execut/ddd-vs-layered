package domain

import (
    "errors"
    "fmt"
)

var (
    ErrLabelTemplateAlreadyCreated = errors.New("попытка создать уже существующий шаблон")
    ErrLabelTemplateAlreadyDeleted = errors.New("попытка удалить уже удалённый шаблон")
    ErrCategoryAlreadyAdded        = errors.New("категория уже привязана к шаблону")
)

type LabelTemplate struct {
    Status       LabelTemplateStatus
    Manufacturer Manufacturer
    ID           LabelTemplateID
    Events       []LabelTemplateEvent
    CategoryList []Category
}

func NewLabelTemplate(id LabelTemplateID) (*LabelTemplate, error) {
    return &LabelTemplate{
        ID:     id,
        Status: LabelTemplateStatusDraft,
    }, nil
}

func (t *LabelTemplate) Create(manufacturer Manufacturer) error {
    if t.Status != LabelTemplateStatusDraft && t.Status != LabelTemplateStatusDeleted {
        return ErrLabelTemplateAlreadyCreated
    }

    err := t.addAndApplyEvent(LabelTemplateCreatedEvent{Manufacturer: manufacturer})
    if err != nil {
        return err
    }

    return nil
}

func (t *LabelTemplate) Update(manufacturer Manufacturer) error {
    err := t.addAndApplyEvent(LabelTemplateUpdatedEvent{Manufacturer: manufacturer})
    if err != nil {
        return err
    }

    return nil
}

func (t *LabelTemplate) AddCategoryList(categoryList []Category) error {
    for _, currentCategory := range t.CategoryList {
        for _, newCategory := range categoryList {
            if currentCategory.Same(newCategory) {
                if newCategory.TypeID != nil {
                    return fmt.Errorf("%w (категория %d, тип %d)", ErrCategoryAlreadyAdded,
                        newCategory.CategoryID, *newCategory.TypeID)
                }

                return fmt.Errorf("%w (категория %d)", ErrCategoryAlreadyAdded, newCategory.CategoryID)
            }
        }
    }

    err := t.addAndApplyEvent(LabelTemplateCategoryListAddedEvent{CategoryList: categoryList})
    if err != nil {
        return err
    }

    return nil
}

func (t *LabelTemplate) Delete() error {
    if t.Status != LabelTemplateStatusCreated {
        return ErrLabelTemplateAlreadyDeleted
    }

    err := t.addAndApplyEvent(LabelTemplateDeletedEvent{})
    if err != nil {
        return err
    }

    return nil
}

func (t *LabelTemplate) ApplyEvent(event LabelTemplateEvent) error {
    switch payload := event.(type) {
    case LabelTemplateCreatedEvent:
        t.Status = LabelTemplateStatusCreated
        t.Manufacturer = payload.Manufacturer
    case LabelTemplateUpdatedEvent:
        t.Manufacturer = payload.Manufacturer
    case LabelTemplateDeletedEvent:
        t.Status = LabelTemplateStatusDeleted
    case LabelTemplateCategoryListAddedEvent:
        t.CategoryList = append(t.CategoryList, payload.CategoryList...)
    }

    return nil
}

func (t *LabelTemplate) addAndApplyEvent(event LabelTemplateEvent) error {
    t.Events = append(t.Events, event)

    err := t.ApplyEvent(event)
    if err != nil {
        return err
    }

    return nil
}
