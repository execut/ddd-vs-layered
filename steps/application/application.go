package application

import (
    "context"
    "encoding/json"

    "effective-architecture/steps/domain"
    "effective-architecture/steps/infrastructure"
)

type Application struct {
    repository domain.IRepository
}

func NewApplication(repository *infrastructure.EventsRepository) (*Application, error) {
    return &Application{
        repository: infrastructure.NewRepository(repository),
    }, nil
}

func (a *Application) CreateLabelTemplate(ctx context.Context, id string, manufacturerOrganizationName string) error {
    domainLabel, err := a.loadLabelTemplate(ctx, id)
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

func (a *Application) GetLabelTemplate(ctx context.Context, id string) (string, error) {
    domainLabel, err := a.loadLabelTemplate(ctx, id)
    if err != nil {
        return "", err
    }

    responseObj := GetLabelTemplateResponse{
        ID:                           domainLabel.ID.UUID.String(),
        ManufacturerOrganizationName: domainLabel.ManufacturerOrganizationName.Name,
    }

    responseMarshaled, err := json.Marshal(responseObj)
    if err != nil {
        return "", err
    }

    return string(responseMarshaled), nil
}

func (a *Application) DeleteLabelTemplate(ctx context.Context, id string) error {
    domainLabel, err := a.loadLabelTemplate(ctx, id)
    if err != nil {
        return err
    }

    err = domainLabel.Delete()
    if err != nil {
        return err
    }

    err = a.repository.Save(ctx, domainLabel)
    if err != nil {
        return err
    }

    return nil
}

func (a *Application) UpdateLabelTemplate(ctx context.Context, uuid string, manufacturerOrganizationName string) error {
    domainLabel, err := a.loadLabelTemplate(ctx, uuid)
    if err != nil {
        return err
    }

    newDomainManufacturerOrganizationName, err := domain.NewManufacturerOrganizationName(manufacturerOrganizationName)
    if err != nil {
        return err
    }

    err = domainLabel.Update(newDomainManufacturerOrganizationName)
    if err != nil {
        return err
    }

    err = a.repository.Save(ctx, domainLabel)
    if err != nil {
        return err
    }

    return nil
}

func (a *Application) loadLabelTemplate(ctx context.Context, uuid string) (*domain.LabelTemplate, error) {
    domainUUID, err := domain.NewLabelTemplateID(uuid)
    if err != nil {
        return nil, err
    }

    domainLabel, err := domain.NewLabelTemplate(domainUUID)
    if err != nil {
        return nil, err
    }

    err = a.repository.Load(ctx, domainLabel)
    if err != nil {
        return nil, err
    }

    return domainLabel, nil
}

type GetLabelTemplateResponse struct {
    ID                           string `json:"id"`
    ManufacturerOrganizationName string `json:"manufacturerOrganizationName"`
}
