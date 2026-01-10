package domain

import "errors"

var (
    ErrLabelTemplateAlreadyCreated = errors.New("попытка создать уже существующий шаблон")
    ErrLabelTemplateAlreadyDeleted = errors.New("попытка удалить уже удалённый шаблон")
)

type LabelTemplate struct {
    Status                       LabelTemplateStatus
    ManufacturerOrganizationName ManufacturerOrganizationName
    ID                           LabelTemplateID
    Events                       []LabelTemplateEvent
}

func NewLabelTemplate(id LabelTemplateID) (*LabelTemplate, error) {
    return &LabelTemplate{
        ID:     id,
        Status: LabelTemplateStatusDraft,
    }, nil
}

func (t *LabelTemplate) Create(manufacturerOrganizationName ManufacturerOrganizationName) error {
    if t.Status != LabelTemplateStatusDraft {
        return ErrLabelTemplateAlreadyCreated
    }

    err := t.addAndApplyEvent(LabelTemplateCreatedEvent{ManufacturerOrganizationName: manufacturerOrganizationName})
    if err != nil {
        return err
    }

    return nil
}

func (t *LabelTemplate) Update(manufacturerOrganizationName ManufacturerOrganizationName) error {
    err := t.addAndApplyEvent(LabelTemplateUpdatedEvent{ManufacturerOrganizationName: manufacturerOrganizationName})
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
        t.ManufacturerOrganizationName = payload.ManufacturerOrganizationName
    case LabelTemplateUpdatedEvent:
        t.ManufacturerOrganizationName = payload.ManufacturerOrganizationName
    case LabelTemplateDeletedEvent:
        t.Status = LabelTemplateStatusDeleted
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
