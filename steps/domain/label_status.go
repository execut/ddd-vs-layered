package domain

type LabelStatus string

const (
    LabelStatusDraft      LabelStatus = "Черновик"
    LabelStatusGeneration LabelStatus = "Генерируется"
)
