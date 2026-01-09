package application

import (
    "context"

    domain2 "effective-architecture/steps/domain"
    infrastructure2 "effective-architecture/steps/infrastructure"
)

type Application struct {
    repository domain2.IRepository
}

func NewApplication(repository *infrastructure2.LabelTemplateRepository) (*Application, error) {
    return &Application{
        repository: infrastructure2.NewRepository(repository),
    }, nil
}

func (a *Application) CreateLabelTemplate(ctx context.Context, uuid string, manufacturerOrganizationName string) error {
    domainUUID, err := domain2.NewLabelTemplateID(uuid)
    if err != nil {
        return err
    }

    domainLabel, err := domain2.NewLabelTemplate(domainUUID)
    if err != nil {
        return err
    }

    err = a.repository.Load(ctx, domainLabel)
    if err != nil {
        return err
    }

    domainManufacturerOrganizationName, err := domain2.NewManufacturerOrganizationName(manufacturerOrganizationName)
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
