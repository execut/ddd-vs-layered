package application

import (
    "context"
    "errors"
)

var ErrCategoryEmpty = errors.New("wrong category")

type Application struct {
}

func NewApplication() (*Application, error) {
    return &Application{}, nil
}

func (a *Application) Create(ctx context.Context, labelTemplateID string, manufacturer Manufacturer) error {
    _ = ctx
    _ = labelTemplateID
    _ = manufacturer

    return nil
}

func (a *Application) Get(ctx context.Context, labelTemplateID string) (LabelTemplate, error) {
    _ = ctx
    _ = labelTemplateID

    return LabelTemplate{}, nil
}

func (a *Application) Delete(ctx context.Context, labelTemplateID string) error {
    _ = ctx
    _ = labelTemplateID

    return nil
}

func (a *Application) Update(ctx context.Context, labelTemplateID string, manufacturer Manufacturer) error {
    _ = ctx
    _ = labelTemplateID
    _ = manufacturer

    return nil
}

func (a *Application) HistoryList(ctx context.Context, labelTemplateID string) ([]LabelTemplateHistoryRow, error) {
    _ = ctx
    _ = labelTemplateID

    return []LabelTemplateHistoryRow{}, nil
}

func (a *Application) AddCategoryList(ctx context.Context, labelTemplateID string, categoryList []Category) error {
    _ = ctx
    _ = labelTemplateID
    _ = categoryList

    return nil
}
