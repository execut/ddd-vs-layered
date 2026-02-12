package service

import (
    "context"
    "errors"
    "strings"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

var (
    ErrLabelAlreadyExists  = errors.New("генерация этикетки с таким идентификатором уже существует")
    ErrLabelCouldNotCreate = errors.New("could not create label")
    ErrLabelCouldNotDelete = errors.New("could not delete label")
)

type LabelRepository struct {
    conn *pgx.Conn
}

func NewLabelRepository(ctx context.Context) (*LabelRepository, error) {
    conn, err := NewDBConnection(ctx)
    if err != nil {
        return nil, err
    }

    return &LabelRepository{conn: conn}, nil
}

func (r LabelRepository) Exists(ctx context.Context, id string) error {
    sql := `SELECT 1 FROM label WHERE id = $1`

    var result int

    err := r.conn.QueryRow(ctx, sql, id).Scan(&result)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil
        }

        return err
    }

    return ErrLabelAlreadyExists
}

func (r LabelRepository) Create(ctx context.Context, model Label) error {
    sql := `
       INSERT INTO label (id, label_template_id, sku)
       VALUES ($1, $2, $3)
   `

    result, err := r.conn.Exec(ctx, sql, model.ID, model.LabelTemplateID, model.SKU)
    if err != nil {
        const duplicateKeyErr = "duplicate key value violates unique constraint"
        if strings.Contains(err.Error(), duplicateKeyErr) {
            return ErrLabelAlreadyExists
        }

        return err
    }

    if result.RowsAffected() == 0 {
        return ErrLabelCouldNotCreate
    }

    return nil
}

func (r LabelRepository) Delete(ctx context.Context, labelID string) error {
    var (
        result pgconn.CommandTag
        err    error
    )

    sql := `DELETE FROM label WHERE id = $1`

    result, err = r.conn.Exec(ctx, sql, labelID)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrLabelCouldNotDelete
    }

    return nil
}
