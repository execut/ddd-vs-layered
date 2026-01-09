package application

import (
    "context"

    "effective-architecture/steps/2-ddd-event-sourcing/domain"
    "effective-architecture/steps/2-ddd-event-sourcing/infrastructure"
)

type Application struct {
    repository domain.IRepository
}

func NewApplication(repository *infrastructure.LabelTemplateRepository) (*Application, error) {
    return &Application{
        repository: infrastructure.NewRepository(repository),
    }, nil
}

func (a *Application) CreateLabelTemplate(ctx context.Context, uuid string, manufacturerOrganizationName string) error {
    domainUUID, err := domain.NewLabelTemplateID(uuid)
    if err != nil {
        return err
    }

    domainLabel, err := domain.NewLabelTemplate(domainUUID)
    if err != nil {
        return err
    }

    err = a.repository.Load(ctx, domainLabel)
    if err != nil {
        return err
    }

    domainManufacturerOrganizationName, err := domain.NewManufacturerOrganizationName(manufacturerOrganizationName)
    if err != nil {
        return err
    }

    err = domainLabel.Create(domainManufacturerOrganizationName)
    if err != nil {
        return err
    }

    err = a.repository.Save(ctx, domainLabel)
    if err != nil {
        return err
    }

    return nil
}
