package infrastructure

import "time"

type EventModel struct {
    AggregateID string
    Type        string
    Payload     []byte
    CreatedAt   time.Time
}
