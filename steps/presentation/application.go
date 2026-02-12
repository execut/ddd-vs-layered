package presentation

import (
    "context"

    "effective-architecture/steps/contract"
    "effective-architecture/steps/contract/external"
    "effective-architecture/steps/service"
)

var _ contract.IApplication = (*Application)(nil)

type Application struct {
}

func NewApplication(ctx context.Context,
    externalServiceOzon external.IExternalServiceOzon) (*Application, error) {
    _ = externalServiceOzon

    repository, err := service.NewRepository(ctx)
    if err != nil {
        return nil, err
    }

    historyRepository, err := service.NewHistoryRepository(ctx)
    if err != nil {
        return nil, err
    }

    categoryRepository, err := service.NewVsCategoryRepository(ctx)
    if err != nil {
        return nil, err
    }

    return service.NewService(repository, historyRepository, categoryRepository), nil
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

func (a *Application) UnlinkCategoryList(ctx context.Context, labelTemplateID string,
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

func (a *Application) IDByCategoryWithType(ctx context.Context,
    categoryWithType contract.CategoryWithType) (string, error) {
    _ = ctx
    _ = categoryWithType

    return "", nil
}

func (a *Application) StartLabelGeneration(ctx context.Context, labelTemplateID string, sku int64) error {
    _ = ctx
    _ = labelTemplateID
    _ = sku

    return nil
}
