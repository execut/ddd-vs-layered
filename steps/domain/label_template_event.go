package domain

type LabelTemplateEvent any

type LabelTemplateCreatedEvent struct {
	Manufacturer Manufacturer
	UserID       string
}

type LabelTemplateUpdatedEvent struct {
	Manufacturer Manufacturer
}

type LabelTemplateDeletedEvent struct {
}

type LabelTemplateActivatedEvent struct {
}

type LabelTemplateDeactivatedEvent struct {
}

type LabelTemplateCategoryListAddedEvent struct {
	CategoryList []Category
}

type LabelTemplateCategoryListUnlinkedEvent struct {
	CategoryList []Category
}
