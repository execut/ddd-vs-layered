package presentation

import (
    "context"

    "effective-architecture/steps/contract"
)

var _ contract.IApplication = (*Application)(nil)

type Application struct {
}

func NewApplication(ctx context.Context) (*Application, error) {
    _ = ctx

    return &Application{}, nil
}

func (a *Application) Create(ctx context.Context, labelTemplateID string, manufacturer contract.Manufacturer) error {
    _ = ctx
    _ = labelTemplateID
    _ = manufacturer

    return nil
}

func (a *Application) Get(ctx context.Context, labelTemplateID string) (contract.LabelTemplate, error) {
    _ = ctx
    _ = labelTemplateID

    return contract.LabelTemplate{}, nil
}

func (a *Application) Delete(ctx context.Context, labelTemplateID string) error {
    _ = ctx
    _ = labelTemplateID

    return nil
}

func (a *Application) Update(ctx context.Context, labelTemplateID string, manufacturer contract.Manufacturer) error {
    _ = ctx
    _ = labelTemplateID
    _ = manufacturer

    return nil
}

func (a *Application) HistoryList(
    ctx context.Context,
    labelTemplateID string) ([]contract.LabelTemplateHistoryRow, error) {
    _ = ctx
    _ = labelTemplateID

    return []contract.LabelTemplateHistoryRow{}, nil
}

func (a *Application) AddCategoryList(ctx context.Context, labelTemplateID string,
    categoryList []contract.Category) error {
    _ = ctx
    _ = labelTemplateID
    _ = categoryList

    return nil
}

func (a *Application) Cleanup(ctx context.Context, labelTemplateID string) error {
    _ = ctx
    _ = labelTemplateID

    return nil
}
