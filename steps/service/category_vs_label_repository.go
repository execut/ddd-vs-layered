package service

import (
    "context"
    "errors"
    "strings"

    "effective-architecture/steps/contract/external"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

var (
    ErrTemplateIDNotFound            = errors.New("шаблон этикетки для SKU не найден")
    ErrCategoryVsLabelAlreadyDeleted = errors.New("category vs label already deleted")
    ErrCategoryVsLabelAlreadyExists  = errors.New("record already exists")
)

type CategoryVsLabelTemplateRepository struct {
    conn *pgx.Conn
}

func NewCategoryVsLabelTemplateRepository(ctx context.Context) (*CategoryVsLabelTemplateRepository, error) {
    conn, err := NewDBConnection(ctx)
    if err != nil {
        return nil, err
    }

    return &CategoryVsLabelTemplateRepository{
        conn: conn,
    }, nil
}

func (r *CategoryVsLabelTemplateRepository) LabelTemplateID(ctx context.Context,
    product external.Product) (string, error) {
    sql := `SELECT label_template_id FROM category_id_vs_label_template_id WHERE category_id = $1 AND type_id=$2`

    var labelTemplateID string

    err := r.conn.QueryRow(ctx, sql, product.Category.Category.ID, product.Category.TypeID).Scan(&labelTemplateID)
    if err != nil && !errors.Is(err, pgx.ErrNoRows) {
        return "", err
    }

    if labelTemplateID != "" {
        return labelTemplateID, nil
    }

    category := product.Category.Category

    for {
        sql = `SELECT label_template_id FROM category_id_vs_label_template_id WHERE category_id = $1
                                                                 AND type_id is null`

        err = r.conn.QueryRow(ctx, sql, category.ID).Scan(&labelTemplateID)
        if err != nil && !errors.Is(err, pgx.ErrNoRows) {
            return "", err
        }

        if labelTemplateID != "" {
            return labelTemplateID, nil
        }

        if category.ParentCategory == nil {
            return "", ErrTemplateIDNotFound
        }

        category = *category.ParentCategory
    }
}

func (r *CategoryVsLabelTemplateRepository) Create(ctx context.Context, model CategoryIDVsLabelTemplateID) error {
    sql := `
       INSERT INTO category_id_vs_label_template_id (label_template_id, category_id, type_id, created_at)
       VALUES ($1, $2, $3, now())
   `

    result, err := r.conn.Exec(ctx, sql, model.LabelTemplateID, model.CategoryID, model.TypeID)
    if err != nil {
        const duplicateKeyErr = "duplicate key value violates unique constraint \"label_id_pk\""
        if strings.Contains(err.Error(), duplicateKeyErr) {
            return ErrCategoryVsLabelAlreadyExists
        }

        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    return nil
}

func (r *CategoryVsLabelTemplateRepository) Delete(ctx context.Context, model CategoryIDVsLabelTemplateID) error {
    var (
        result pgconn.CommandTag
        err    error
    )

    sql := `DELETE FROM category_id_vs_label_template_id WHERE label_template_id = $1 AND category_id=$2`
    if model.TypeID == nil {
        result, err = r.conn.Exec(ctx, sql, model.LabelTemplateID, model.CategoryID)
        if err != nil {
            return err
        }
    } else {
        result, err = r.conn.Exec(ctx, sql+" AND type_id=$3", model.LabelTemplateID, model.CategoryID, *model.TypeID)
        if err != nil {
            return err
        }
    }

    if result.RowsAffected() == 0 {
        return ErrCategoryVsLabelAlreadyDeleted
    }

    return nil
}

func (r *CategoryVsLabelTemplateRepository) DeleteByLabelTemplateID(ctx context.Context, labelTemplateID string) error {
    sql := `DELETE FROM category_id_vs_label_template_id WHERE label_template_id = $1`

    result, err := r.conn.Exec(ctx, sql, labelTemplateID)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotDelete
    }

    return nil
}
