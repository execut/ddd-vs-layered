package domain

import "github.com/google/uuid"

type LabelID struct {
    UUID uuid.UUID
}

func NewLabelID(id string) (LabelID, error) {
    uuidValue, _ := uuid.Parse(id)

    return LabelID{
        UUID: uuidValue,
    }, nil
}
