package domain

type LabelStatus string

const (
	LabelStatusGeneration LabelStatus = "Генерация"
	LabelStatusDataFilled LabelStatus = "Наполнено данными"
	LabelStatusGenerated  LabelStatus = "Сгенерировано"
)
