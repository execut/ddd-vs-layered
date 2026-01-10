package labels

import (
    "context"
    "errors"
    "os"

    "github.com/jackc/pgx/v5"
)

var (
    ErrCouldNotTruncate = errors.New("could not truncate label template")
    ErrCouldNotCreate   = errors.New("could not create label template")
    ErrCouldNotDelete   = errors.New("could not delete label template")
)

type Repository struct {
    conn *pgx.Conn
}

func NewRepository(ctx context.Context) (*Repository, error) {
    connString := os.Getenv("DATABASE_URL")

    conn, err := pgx.Connect(ctx, connString)
    if err != nil {
        return nil, err
    }

    err = conn.Ping(ctx)
    if err != nil {
        return nil, err
    }

    return &Repository{conn: conn}, nil
}

func (r Repository) Insert(ctx context.Context, model LabelTemplate) error {
    sql := `
        INSERT INTO label_templates (id, manufacturer_organization_name)
        VALUES ($1, $2)
    `

    result, err := r.conn.Exec(ctx, sql, model.ID, model.ManufacturerOrganizationName)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    return nil
}

func (r Repository) Find(ctx context.Context, id string) (LabelTemplate, error) {
    sql := "SELECT id, manufacturer_organization_name FROM label_templates WHERE id = $1"
    model := LabelTemplate{}

    err := r.conn.QueryRow(ctx, sql, id).Scan(&model.ID, &model.ManufacturerOrganizationName)
    if err != nil {
        return LabelTemplate{}, err
    }

    return model, nil
}

func (r Repository) Truncate(ctx context.Context) error {
    sql := `
    TRUNCATE label_templates
    `

    result, err := r.conn.Exec(ctx, sql)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotTruncate
    }

    return nil
}

func (r Repository) Delete(ctx context.Context, id string) error {
    result, err := r.conn.Exec(ctx, "DELETE FROM label_templates WHERE id = $1", id)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotTruncate
    }

    return nil
}
