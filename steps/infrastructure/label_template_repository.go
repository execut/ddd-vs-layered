package infrastructure

import (
    "context"
    "errors"
    "fmt"
    "os"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
)

var (
    ErrCouldNotTruncate = errors.New("could not truncate label template")
    ErrCouldNotCreate   = errors.New("could not create label template")
)

type LabelTemplateRepository struct {
    conn *pgx.Conn
}

func NewLabelTemplateRepository() (*LabelTemplateRepository, error) {
    conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %w", err)
    }

    return &LabelTemplateRepository{
        conn: conn,
    }, nil
}

func (r *LabelTemplateRepository) Create(ctx context.Context, model LabelTemplate) error {
    sql := "INSERT INTO label_templates VALUES ($1, $2)"

    result, err := r.conn.Exec(ctx, sql, model.ID, model.ManufacturerOrganizationName)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    return nil
}

func (r *LabelTemplateRepository) Truncate(ctx context.Context) error {
    result, err := r.conn.Exec(ctx, "TRUNCATE label_templates")
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotTruncate
    }

    return nil
}

func (r *LabelTemplateRepository) Update(_ context.Context, _ LabelTemplate) error {
    return nil
}

func (r *LabelTemplateRepository) Load(ctx context.Context, labelTemplateID uuid.UUID) (LabelTemplate, error) {
    template := LabelTemplate{}

    sql := "SELECT id, manufacturer_organization_name FROM label_templates WHERE id=$1"

    err := r.conn.QueryRow(ctx, sql, labelTemplateID.String()).
        Scan(&template.ID, &template.ManufacturerOrganizationName)
    if err != nil {
        return LabelTemplate{}, fmt.Errorf("error select label template: %w", err)
    }

    return template, nil
}
