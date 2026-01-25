package service

import (
    "context"
    "errors"
    "os"

    "github.com/jackc/pgx/v5"
)

var (
    ErrCouldNotCreate = errors.New("could not create label template")
    ErrCouldNotDelete   = errors.New("could not delete label template")
    ErrCouldNotUpdate   = errors.New("could not update label template")
)

type Repository struct {
    conn *pgx.Conn
}

func NewRepository(ctx context.Context) (*Repository, error) {
    conn, err := NewDBConnection(ctx)
    if err != nil {
        return nil, err
    }

    return &Repository{conn: conn}, nil
}

func NewDBConnection(ctx context.Context) (*pgx.Conn, error) {
    connString := os.Getenv("DATABASE_URL")

    conn, err := pgx.Connect(ctx, connString)
    if err != nil {
        return nil, err
    }

    err = conn.Ping(ctx)
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func (r Repository) Insert(ctx context.Context, model LabelTemplate) error {
    sql := `
        INSERT INTO label_templates (id, manufacturer_organization_name, manufacturer_organization_address,
            manufacturer_email, manufacturer_site)
        VALUES ($1, $2, $3, $4, $5)
    `

    result, err := r.conn.Exec(ctx, sql, model.ID, model.ManufacturerOrganizationName,
        model.ManufacturerOrganizationAddress, model.ManufacturerEmail, model.ManufacturerSite)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    return nil
}

func (r Repository) Update(ctx context.Context, model LabelTemplate) error {
    sql := `
        UPDATE label_templates SET id=$1, manufacturer_organization_name=$2,
            manufacturer_organization_address=$3,
            manufacturer_email=$4,
            manufacturer_site=$5 WHERE id=$1
    `

    result, err := r.conn.Exec(ctx, sql, model.ID, model.ManufacturerOrganizationName,
        model.ManufacturerOrganizationAddress, model.ManufacturerEmail, model.ManufacturerSite)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotUpdate
    }

    return nil
}

func (r Repository) Find(ctx context.Context, id string) (LabelTemplate, error) {
    sql := `SELECT id, manufacturer_organization_name, manufacturer_organization_address,
       manufacturer_email, manufacturer_site FROM label_templates WHERE id = $1`
    model := LabelTemplate{}

    err := r.conn.QueryRow(ctx, sql, id).Scan(&model.ID, &model.ManufacturerOrganizationName,
        &model.ManufacturerOrganizationAddress, &model.ManufacturerEmail, &model.ManufacturerSite)
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
        return ErrCouldNotDelete
    }

    return nil
}

func (r Repository) Delete(ctx context.Context, id string) error {
    result, err := r.conn.Exec(ctx, "DELETE FROM label_templates WHERE id = $1", id)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotDelete
    }

    return nil
}
