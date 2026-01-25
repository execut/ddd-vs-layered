package service

import (
    "context"
    "errors"
    "strconv"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
)

var ErrFailedToCreateHistoryCategory = errors.New("failed to create label_templates_history_categories")

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
    model.ID = uuid.New()
    sql := `
        INSERT INTO label_templates_history (id, label_template_id, order_key, action,
             new_manufacturer_organization_name, new_manufacturer_organization_address, new_manufacturer_email,
                                             new_manufacturer_site, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
    `

    result, err := r.conn.Exec(ctx, sql, model.ID, model.LabelTemplateID, model.OrderKey, model.Action,
        model.NewManufacturerOrganizationName, model.NewManufacturerOrganizationAddress, model.NewManufacturerEmail,
        model.NewManufacturerSite)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    for _, category := range model.CategoryList {
        err = r.createCategory(ctx, model.ID.String(), category)
        if err != nil {
            return err
        }
    }

    return nil
}

func (r HistoryRepository) FindAll(ctx context.Context, labelTemplateID string) ([]LabelTemplateHistory, error) {
    sql := `SELECT id, order_key, action, new_manufacturer_organization_name,
                     new_manufacturer_organization_address, new_manufacturer_email, new_manufacturer_site
            FROM label_templates_history WHERE label_template_id = $1`

    rows, err := r.conn.Query(ctx, sql, labelTemplateID)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var modelList []LabelTemplateHistory

    for rows.Next() {
        var model LabelTemplateHistory

        err := rows.Scan(&model.ID, &model.OrderKey, &model.Action, &model.NewManufacturerOrganizationName,
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

    newModelList := make([]LabelTemplateHistory, 0, len(modelList))
    for _, model := range modelList {
        model.CategoryList, err = r.findCategoryList(ctx, model.ID.String())
        if err != nil {
            return nil, err
        }

        newModelList = append(newModelList, model)
    }

    return newModelList, nil
}

func (r HistoryRepository) Delete(ctx context.Context, labelTemplateID string) error {
    sql := `DELETE FROM label_templates_history WHERE label_template_id = $1`

    result, err := r.conn.Exec(ctx, sql, labelTemplateID)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotDelete
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

func (r HistoryRepository) createCategory(ctx context.Context, historyID string, model HistoryCategory) error {
    sql := `
        INSERT INTO label_templates_history_categories (history_id, category_id, type_id)
        VALUES ($1, $2, $3)
    `

    result, err := r.conn.Exec(ctx, sql, historyID, model.CategoryID, model.TypeID)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrFailedToCreateHistoryCategory
    }

    return nil
}

func (r HistoryRepository) findCategoryList(ctx context.Context, historyID string) ([]HistoryCategory, error) {
    sql := `SELECT category_id, type_id
            FROM label_templates_history_categories WHERE history_id = $1`

    rows, err := r.conn.Query(ctx, sql, historyID)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    var modelList []HistoryCategory

    for rows.Next() {
        var model HistoryCategory

        err := rows.Scan(&model.CategoryID, &model.TypeID)
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
