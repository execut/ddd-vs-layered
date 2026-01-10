package domain

type LabelTemplateEvent any

type LabelTemplateCreatedEvent struct {
    ManufacturerOrganizationName ManufacturerOrganizationName
}

type LabelTemplateUpdatedEvent struct {
    ManufacturerOrganizationName ManufacturerOrganizationName
}

type LabelTemplateDeletedEvent struct {
}
