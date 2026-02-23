package domain

type LabelEvent any

type LabelGenerationStartedEvent struct {
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
