package domain

type LabelEvent any

type LabelGenerationStartedEvent struct {
	UserID          string
	LabelTemplateID LabelTemplateID
	SKU             int64
}

type LabelDataFilledEvent struct {
	LabelTemplateID LabelTemplateID
	Product         Product
	Manufacturer    Manufacturer
}

type LabelGeneratedEvent struct {
	LabelFile LabelFile
}
