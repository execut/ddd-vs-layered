package presentation

import (
    "context"
    "fmt"
    "os"

    "effective-architecture/steps/application"
    "effective-architecture/steps/contract"
    "effective-architecture/steps/infrastructure"
    "effective-architecture/steps/infrastructure/history"

    "github.com/jackc/pgx/v5"
)

var _ contract.IApplication = (*Application)(nil)

type Application struct {
}

func NewApplication(ctx context.Context) (*application.Application, error) {
    conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %w", err)
    }

    eventRepository, err := infrastructure.NewEventsRepository(conn)
    if err != nil {
        return nil, err
    }

    repository := infrastructure.NewRepository(eventRepository)

    historyRepository, err := history.NewRepository(conn)
    if err != nil {
        return nil, err
    }

    app, err := application.NewApplication(repository, historyRepository)
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
