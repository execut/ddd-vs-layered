package domain

type LabelEvent any

type LabelGenerationStartedEvent struct {
    LabelTemplateID LabelTemplateID
    SKU int64
}
