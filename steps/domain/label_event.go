package domain

type LabelEvent any

type LabelGenerationStartedEvent struct {
	LabelTemplateID LabelTemplateID
	SKU             int64
}

type LabelDataFilledEvent struct {
	LabelTemplateID LabelTemplateID
	Product         Product
}
