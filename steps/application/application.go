package application

import (
    "context"
    "encoding/json"

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

func (a *Application) GetLabelTemplate(ctx context.Context, uuid string) (string, error) {
    domainUUID, err := domain2.NewLabelTemplateID(uuid)
    if err != nil {
        return "", err
    }

    domainLabel, err := domain2.NewLabelTemplate(domainUUID)
    if err != nil {
        return "", err
    }

    err = a.repository.Load(ctx, domainLabel)
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
    domainUUID, err := domain2.NewLabelTemplateID(id)
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

type GetLabelTemplateResponse struct {
    ID                           string `json:"id"`
    ManufacturerOrganizationName string `json:"manufacturerOrganizationName"`
}
