package domain

import "github.com/google/uuid"

type LabelTemplateID struct {
    UUID uuid.UUID
}

func NewLabelTemplateID(id string) (LabelTemplateID, error) {
    uuidValue, _ := uuid.Parse(id)

    return LabelTemplateID{
        UUID: uuidValue,
    }, nil
}
