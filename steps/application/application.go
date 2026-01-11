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

func (a *Application) CreateLabelTemplate(ctx context.Context, id string, manufacturer Manufacturer) error {
    domainLabel, err := a.loadLabelTemplate(ctx, id)
    if err != nil {
        return err
    }

    domainManufacturer, err := mapManufacturerToDomain(manufacturer)
    if err != nil {
        return err
    }

    err = domainLabel.Create(domainManufacturer)
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

    response := mapManufacturerToResponse(domainLabel.Manufacturer)
    responseObj := GetLabelTemplateResponse{
        ID:           domainLabel.ID.UUID.String(),
        Manufacturer: response,
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

func (a *Application) UpdateLabelTemplate(ctx context.Context, uuid string, manufacturer Manufacturer) error {
    domainLabel, err := a.loadLabelTemplate(ctx, uuid)
    if err != nil {
        return err
    }

    domainManufacturer, err := mapManufacturerToDomain(manufacturer)
    if err != nil {
        return err
    }

    err = domainLabel.Update(domainManufacturer)
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
    ID           string       `json:"id"`
    Manufacturer Manufacturer `json:"manufacturer"`
}

type Manufacturer struct {
    OrganizationName    string `json:"organizationName"`
    OrganizationAddress string `json:"organizationAddress,omitempty"`
    Site                string `json:"site,omitempty"`
    Email               string `json:"email,omitempty"`
}

func mapManufacturerToDomain(manufacturer Manufacturer) (domain.Manufacturer, error) {
    newDomainManufacturerOrganizationName, err := domain.NewOrganizationName(manufacturer.OrganizationName)
    if err != nil {
        return domain.Manufacturer{}, err
    }

    var (
        domainOrganizationAddress *domain.OrganizationAddress
        domainSite                *domain.Site
        domainEmail               *domain.Email
    )

    if manufacturer.OrganizationAddress != "" {
        domainOrganizationAddressValue, err := domain.NewOrganizationAddress(manufacturer.OrganizationAddress)
        if err != nil {
            return domain.Manufacturer{}, err
        }

        domainOrganizationAddress = &domainOrganizationAddressValue
    }

    if manufacturer.Site != "" {
        domainSiteValue, err := domain.NewSite(manufacturer.Site)
        if err != nil {
            return domain.Manufacturer{}, err
        }

        domainSite = &domainSiteValue
    }

    if manufacturer.Email != "" {
        domainEmailValue, err := domain.NewEmail(manufacturer.Email)
        if err != nil {
            return domain.Manufacturer{}, err
        }

        domainEmail = &domainEmailValue
    }

    domainManufacturer, err := domain.NewManufacturer(newDomainManufacturerOrganizationName, domainOrganizationAddress,
        domainSite, domainEmail)
    if err != nil {
        return domain.Manufacturer{}, err
    }

    return domainManufacturer, nil
}

func mapManufacturerToResponse(domainManufacturer domain.Manufacturer) Manufacturer {
    manufacturer := Manufacturer{
        OrganizationName: domainManufacturer.OrganizationName.Name,
    }

    if domainManufacturer.OrganizationAddress != nil {
        manufacturer.OrganizationAddress = string(*domainManufacturer.OrganizationAddress)
    }

    if domainManufacturer.Site != nil {
        manufacturer.Site = string(*domainManufacturer.Site)
    }

    if domainManufacturer.Email != nil {
        manufacturer.Email = string(*domainManufacturer.Email)
    }

    return manufacturer
}
