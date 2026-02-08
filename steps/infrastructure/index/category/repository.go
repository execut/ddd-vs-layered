package category

import (
    "context"
    "errors"
    "fmt"
    "time"

    "effective-architecture/steps/domain"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgtype"
)

var ErrNotFound = errors.New("category not found")

type Repository struct {
    conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) (*Repository, error) {
    return &Repository{
        conn: conn,
    }, nil
}

func (r Repository) Create(ctx context.Context, index Index) error {
    sql := "INSERT INTO category_vs_label_template VALUES ($1, $2, $3, $4)"

    _, err := r.conn.Exec(ctx, sql, index.CategoryID, index.TypeID,
        index.LabelTemplateID, time.Now())
    if err != nil {
        return err
    }

    return nil
}

func (r Repository) Delete(ctx context.Context, index Index) error {
    sql := "DELETE FROM category_vs_label_template WHERE category_id=$1 AND label_template_id=$2 AND "

    args := []any{index.CategoryID, index.LabelTemplateID}
    if index.TypeID == nil {
        sql += " type_id IS NULL"
    } else {
        args = append(args, *index.TypeID)
        sql += " type_id=$3"
    }

    _, err := r.conn.Exec(ctx, sql, args...)
    if err != nil {
        return err
    }

    return nil
}

func (r Repository) AggregateIDByCategory(ctx context.Context, category domain.Category) (*string, error) {
    sql := "SELECT label_template_id FROM category_vs_label_template WHERE category_id=$1 AND "
    args := []any{category.CategoryID}

    if category.TypeID == nil {
        sql += "type_id IS NULL"
    } else {
        args = append(args, *category.TypeID)
        sql += "type_id=$2"
    }

    var templateID *pgtype.UUID

    err := r.conn.QueryRow(ctx, sql, args...).Scan(&templateID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }

        return nil, fmt.Errorf("failed to count history: %w", err)
    }

    templateIDString := templateID.String()

    return &templateIDString, nil
}

func (r Repository) Cleanup(ctx context.Context, id domain.LabelTemplateID) error {
    _, err := r.conn.Exec(ctx, "DELETE FROM category_vs_label_template WHERE label_template_id=$1", id.UUID)
    if err != nil {
        return err
    }

    return nil
}
