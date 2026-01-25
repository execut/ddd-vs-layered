package service

import (
    "context"
    "errors"
    "fmt"
    "strings"

    "github.com/jackc/pgx/v5"
)

var ErrVsCategoryAlreadyCreated = errors.New("категория уже привязана к шаблону")

type VsCategoryRepository struct {
    conn *pgx.Conn
}

func NewVsCategoryRepository(ctx context.Context) (*VsCategoryRepository, error) {
    conn, err := NewDBConnection(ctx)
    if err != nil {
        return nil, err
    }

    return &VsCategoryRepository{conn: conn}, nil
}

func (r VsCategoryRepository) Create(ctx context.Context, model LabelTemplateVsCategory) error {
    sql := `
        INSERT INTO label_template_vs_categories (label_template_id, category_id, type_id)
        VALUES ($1, $2, $3)
    `

    result, err := r.conn.Exec(ctx, sql, model.LabelTemplateID, model.CategoryID, model.TypeID)
    if err != nil {
        const duplicateKeyErr = "duplicate key value violates unique constraint \"label_template_vs_categories_id_pk\""
        if strings.Contains(err.Error(), duplicateKeyErr) {
            if model.TypeID == nil {
                return fmt.Errorf("%w (категория %v)", ErrVsCategoryAlreadyCreated, model.CategoryID)
            }

            return fmt.Errorf("%w (категория %v, тип %v)", ErrVsCategoryAlreadyCreated, model.CategoryID, *model.TypeID)
        }

        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotCreate
    }

    return nil
}

func (r VsCategoryRepository) Delete(ctx context.Context, labelTemplateID string) error {
    sql := `DELETE FROM label_template_vs_categories WHERE label_template_id = $1`

    result, err := r.conn.Exec(ctx, sql, labelTemplateID)
    if err != nil {
        return err
    }

    if result.RowsAffected() == 0 {
        return ErrCouldNotDelete
    }

    return nil
}
