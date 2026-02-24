package domain

type LabelTemplateStatus string

const (
	LabelTemplateStatusDraft       LabelTemplateStatus = "Черновик"
	LabelTemplateStatusCreated     LabelTemplateStatus = "Создан"
	LabelTemplateStatusDeleted     LabelTemplateStatus = "Удалён"
	LabelTemplateStatusDeactivated LabelTemplateStatus = "Деактивирован"
)
