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
	ErrLabelCouldNotUpdate = errors.New("could not update label")
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

func (r LabelRepository) Get(ctx context.Context, generationID string) (Label, error) {
	sql := `SELECT label_template_id, sku, status, product_name, manufacturer_organization_name,
            manufacturer_organization_address, manufacturer_email, manufacturer_site,
            file FROM label WHERE id = $1`
	model := Label{}

	err := r.conn.QueryRow(ctx, sql, generationID).Scan(&model.LabelTemplateID, &model.SKU,
		&model.Status, &model.ProductName, &model.ManufacturerOrganizationName, &model.ManufacturerOrganizationAddress,
		&model.ManufacturerEmail, &model.ManufacturerSite, &model.File)
	if err != nil {
		return Label{}, err
	}

	model.ID = generationID

	return model, nil
}

func (r LabelRepository) Create(ctx context.Context, model Label) error {
	sql := `
       INSERT INTO label (id, sku, status)
       VALUES ($1, $2, $3)
   `

	result, err := r.conn.Exec(ctx, sql, model.ID, model.SKU, model.Status)
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

func (r LabelRepository) Update(ctx context.Context, model Label) error {
	sql := `
       UPDATE label SET label_template_id=$2, sku=$3, status=$4, product_name=$5, manufacturer_organization_name=$6,
            manufacturer_organization_address=$7,
            manufacturer_email=$8,
            manufacturer_site=$9,
            file=$10
       WHERE id=$1
   `

	result, err := r.conn.Exec(ctx, sql, model.ID, model.LabelTemplateID, model.SKU, model.Status, model.ProductName,
		model.ManufacturerOrganizationName, model.ManufacturerOrganizationAddress, model.ManufacturerEmail,
		model.ManufacturerSite, model.File)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrLabelCouldNotUpdate
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
