package service

import (
    "context"
    "strconv"

    "github.com/jackc/pgx/v5"
)

type HistoryRepository struct {
    conn *pgx.Conn
}

func NewHistoryRepository(ctx context.Context) (*HistoryRepository, error) {
    conn, err := NewDBConnection(ctx)
    if err != nil {
        return nil, err
    }

    return &HistoryRepository{conn: conn}, nil
}

func (r HistoryRepository) Create(ctx context.Context, model LabelTemplateHistory, orderKey int) error {
    if orderKey == 0 {
        var err error

        orderKey, err = r.lastKey(ctx, model.LabelTemplateID)
        if err != nil {
            return err
        }
    }

    model.OrderKey = strconv.Itoa(orderKey + 1)
    sql := `
        INSERT INTO label_templates_history (label_template_id, order_key, action, new_manufacturer_organization_name,
                     new_manufacturer_organization_address, new_manufacturer_email, new_manufacturer_site, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
    `

    result, err := r.conn.Exec(ctx, sql, model.LabelTemplateID, model.OrderKey, model.Action,
        model.NewManufacturerOrganizationName, model.NewManufacturerOrganizationAddress, model.NewManufacturerEmail,
        model.NewManufacturerSite)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    return nil
}

func (r HistoryRepository) FindAll(ctx context.Context, labelTemplateID string) ([]LabelTemplateHistoryResult, error) {
    sql := `SELECT order_key, action, new_manufacturer_organization_name,
                     new_manufacturer_organization_address, new_manufacturer_email, new_manufacturer_site
            FROM label_templates_history WHERE label_template_id = $1`

    rows, err := r.conn.Query(ctx, sql, labelTemplateID)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var modelList []LabelTemplateHistoryResult

    for rows.Next() {
        var model LabelTemplateHistoryResult

        err := rows.Scan(&model.OrderKey, &model.Action, &model.NewManufacturerOrganizationName,
            &model.NewManufacturerOrganizationAddress, &model.NewManufacturerEmail, &model.NewManufacturerSite)
        if err != nil {
            return nil, err
        }

        modelList = append(modelList, model)
    }

    err = rows.Err()
    if err != nil {
        return nil, err
    }

    return modelList, nil
}

func (r HistoryRepository) Truncate(ctx context.Context) error {
    sql := `TRUNCATE label_templates_history`

    result, err := r.conn.Exec(ctx, sql)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotTruncate
    }

    return nil
}

func (r HistoryRepository) lastKey(ctx context.Context, labelTemplateID string) (int, error) {
    sql := `SELECT max(order_key) FROM label_templates_history WHERE label_template_id = $1`

    var lastKey *int

    err := r.conn.QueryRow(ctx, sql, labelTemplateID).Scan(&lastKey)
    if err != nil {
        return 0, err
    }

    if lastKey == nil {
        return 0, nil
    }

    return *lastKey, nil
}
