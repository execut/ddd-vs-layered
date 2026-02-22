package presentation

import (
    "context"
    "fmt"
    "os"

    "effective-architecture/steps/application"
    "effective-architecture/steps/contract"
    "effective-architecture/steps/contract/external"
    "effective-architecture/steps/domain"
    "effective-architecture/steps/infrastructure"
    "effective-architecture/steps/infrastructure/history"
    "effective-architecture/steps/infrastructure/index/category"

    "github.com/jackc/pgx/v5"
)

var _ contract.IApplication = (*Application)(nil)

type Application struct {
}

func NewApplication(ctx context.Context,
    externalServiceOzon external.IExternalServiceOzon) (*application.Application, error) {
    conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %w", err)
    }

    eventRepository, err := infrastructure.NewEventsRepository(conn)
    if err != nil {
        return nil, err
    }

    categoryVsTemplateRepository, err := category.NewRepository(conn)
    if err != nil {
        return nil, err
    }

    repository := infrastructure.NewRepository(eventRepository, categoryVsTemplateRepository)

    historyRepository, err := history.NewRepository(conn)
    if err != nil {
        return nil, err
    }

    ozonService := infrastructure.NewServiceOzon(externalServiceOzon)

    labelRepository := infrastructure.NewLabelRepository(eventRepository)

    app, err := application.NewApplication(repository, historyRepository, []domain.Subscriber{
        category.NewSubscriber(categoryVsTemplateRepository),
    }, labelRepository, ozonService)
    if err != nil {
        return nil, err
    }

    return app, nil
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

func (a *Application) StartLabelGeneration(ctx context.Context, generationID string, sku int64) error {
	_ = ctx
	_ = generationID
	_ = sku

	return nil
}

func (a *Application) LabelGeneration(ctx context.Context, generationID string) (contract.LabelGeneration, error) {
	_ = ctx
	_ = generationID

	return contract.LabelGeneration{}, nil
}
