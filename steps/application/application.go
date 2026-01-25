package application

import (
    "context"

    "effective-architecture/steps/contract"
    "effective-architecture/steps/domain"
    "effective-architecture/steps/domain/history"
)

var _ contract.IApplication = (*Application)(nil)

type Application struct {
    repository        domain.IRepository
    dispatcher        *domain.Dispatcher
    historyRepository history.IRepository
}

func NewApplication(repository domain.IRepository,
    historyRepository history.IRepository) (*Application, error) {
    dispatcher := domain.NewDispatcher([]domain.Subscriber{
        history.NewSubscriber(historyRepository),
    })

    return &Application{
        repository:        repository,
        historyRepository: historyRepository,
        dispatcher:        dispatcher,
    }, nil
}

func (a *Application) Create(ctx context.Context, id string, manufacturer contract.Manufacturer) error {
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

    err = a.dispatcher.Dispatch(ctx, domainLabel)
    if err != nil {
        return err
    }

    err = a.repository.Save(ctx, domainLabel)
    if err != nil {
        return err
    }

    return nil
}

func (a *Application) Get(ctx context.Context, id string) (contract.LabelTemplate, error) {
    domainLabel, err := a.loadLabelTemplate(ctx, id)
    if err != nil {
        return contract.LabelTemplate{}, err
    }

    response := mapManufacturerToResponse(domainLabel.Manufacturer)
    responseObj := contract.LabelTemplate{
        ID:           domainLabel.ID.UUID.String(),
        Manufacturer: response,
    }

    return responseObj, nil
}

func (a *Application) Delete(ctx context.Context, id string) error {
    domainLabel, err := a.loadLabelTemplate(ctx, id)
    if err != nil {
        return err
    }

    err = domainLabel.Delete()
    if err != nil {
        return err
    }

    err = a.dispatcher.Dispatch(ctx, domainLabel)
    if err != nil {
        return err
    }

    err = a.repository.Save(ctx, domainLabel)
    if err != nil {
        return err
    }

    return nil
}

func (a *Application) Cleanup(ctx context.Context, id string) error {
    return nil
}

func (a *Application) Update(ctx context.Context, uuid string, manufacturer contract.Manufacturer) error {
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

    err = a.dispatcher.Dispatch(ctx, domainLabel)
    if err != nil {
        return err
    }

    err = a.repository.Save(ctx, domainLabel)
    if err != nil {
        return err
    }

    return nil
}

func (a *Application) HistoryList(ctx context.Context, id string) ([]contract.LabelTemplateHistoryRow, error) {
    domainAggregateID, err := domain.NewLabelTemplateID(id)
    if err != nil {
        return nil, err
    }

    domainHistoryList, err := a.historyRepository.List(ctx, domainAggregateID)
    if err != nil {
        return nil, err
    }

    result := []contract.LabelTemplateHistoryRow{}

    for _, domainHistory := range domainHistoryList {
        historyRow := contract.LabelTemplateHistoryRow{
            OrderKey: domainHistory.OrderKey,
            Action:   domainHistory.Action,
        }

        if domainHistory.NewManufacturerOrganizationName != nil {
            historyRow.NewManufacturerOrganizationName = domainHistory.NewManufacturerOrganizationName.Name
        }

        if domainHistory.NewManufacturerOrganizationAddress != nil {
            historyRow.NewManufacturerOrganizationAddress = domainHistory.NewManufacturerOrganizationAddress.Address
        }

        if domainHistory.NewManufacturerEmail != nil {
            historyRow.NewManufacturerEmail = domainHistory.NewManufacturerEmail.Value
        }

        if domainHistory.NewManufacturerSite != nil {
            historyRow.NewManufacturerSite = domainHistory.NewManufacturerSite.Value
        }

        result = append(result, historyRow)
    }

    return result, nil
}

func (a *Application) loadLabelTemplate(ctx context.Context, id string) (*domain.LabelTemplate, error) {
    domainID, err := domain.NewLabelTemplateID(id)
    if err != nil {
        return nil, err
    }

    domainLabel, err := domain.NewLabelTemplate(domainID)
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

func mapManufacturerToDomain(manufacturer contract.Manufacturer) (domain.Manufacturer, error) {
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

func mapManufacturerToResponse(domainManufacturer domain.Manufacturer) contract.Manufacturer {
    manufacturer := contract.Manufacturer{
        OrganizationName: domainManufacturer.OrganizationName.Name,
    }

    if domainManufacturer.OrganizationAddress != nil {
        manufacturer.OrganizationAddress = domainManufacturer.OrganizationAddress.Address
    }

    if domainManufacturer.Site != nil {
        manufacturer.Site = domainManufacturer.Site.Value
    }

    if domainManufacturer.Email != nil {
        manufacturer.Email = domainManufacturer.Email.Value
    }

    return manufacturer
}
