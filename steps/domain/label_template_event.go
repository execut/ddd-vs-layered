package domain

type LabelTemplateEvent any

type LabelTemplateCreatedEvent struct {
    Manufacturer Manufacturer
}

type LabelTemplateUpdatedEvent struct {
    Manufacturer Manufacturer
}

type LabelTemplateDeletedEvent struct {
}

type LabelTemplateCategoryListAddedEvent struct {
    CategoryList []Category
}

type LabelTemplateCategoryListUnlinkedEvent struct {
    CategoryList []Category
}
