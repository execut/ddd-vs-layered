package infrastructure

import (
    "context"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
)

var (
    ErrCouldNotCreate = errors.New("could not create label template")
)

type EventsRepository struct {
    conn *pgx.Conn
}

func NewEventsRepository(conn *pgx.Conn) (*EventsRepository, error) {
    return &EventsRepository{
        conn: conn,
    }, nil
}

func (r *EventsRepository) Save(ctx context.Context, modelList []EventModel) error {
    for _, model := range modelList {
        sql := "INSERT INTO label_templates_events VALUES ($1, $2, $3, $4)"

        result, err := r.conn.Exec(ctx, sql, model.AggregateID, model.Type, model.Payload, model.CreatedAt)
        if err != nil {
            return err
        }

        if result.RowsAffected() == 0 {
            return ErrCouldNotCreate
        }
    }

    return nil
}

func (r *EventsRepository) Cleanup(ctx context.Context, aggregateID uuid.UUID) error {
    _, err := r.conn.Exec(ctx, "DELETE FROM label_templates_events WHERE aggregate_id=$1", aggregateID)
    if err != nil {
        return err
    }

    return nil
}

func (r *EventsRepository) Load(ctx context.Context, aggregateID uuid.UUID) ([]EventModel, error) {
    sql := "SELECT aggregate_id, type, payload, created_at FROM label_templates_events WHERE aggregate_id=$1"

    rows, err := r.conn.Query(ctx, sql, aggregateID.String())
    if err != nil {
        return nil, fmt.Errorf("error select label template event: %w", err)
    }

    result := make([]EventModel, 0)

    for rows.Next() {
        var model EventModel

        err := rows.Scan(&model.AggregateID, &model.Type, &model.Payload, &model.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("error scan label template event: %w", err)
        }

        result = append(result, model)
    }

    return result, nil
}
