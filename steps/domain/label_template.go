package domain

type LabelTemplate struct {
    ManufacturerOrganizationName ManufacturerOrganizationName
    ID                           LabelTemplateID
    Events                       []LabelTemplateEvent
}

func NewLabelTemplate(id LabelTemplateID) (*LabelTemplate, error) {
    return &LabelTemplate{
        ID: id,
    }, nil
}

func (t *LabelTemplate) Create(manufacturerOrganizationName ManufacturerOrganizationName) error {
    err := t.addAndApplyEvent(LabelTemplateCreatedEvent{ManufacturerOrganizationName: manufacturerOrganizationName})
    if err != nil {
        return err
    }

    return nil
}

func (t *LabelTemplate) Delete() error {
    err := t.addAndApplyEvent(LabelTemplateDeletedEvent{})
    if err != nil {
        return err
    }

    return nil
}

func (t *LabelTemplate) addAndApplyEvent(event LabelTemplateEvent) error {
    t.Events = append(t.Events, event)
    if payload, ok := event.(LabelTemplateCreatedEvent); ok {
        t.ManufacturerOrganizationName = payload.ManufacturerOrganizationName
    }

    return nil
}
