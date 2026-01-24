package application

import (
    "context"
    "encoding/json"
    "errors"
    "strconv"
    "strings"

    "effective-architecture/steps/domain"
    "effective-architecture/steps/domain/history"
    "effective-architecture/steps/infrastructure"
)

var ErrCategoryEmpty = errors.New("wrong category")

type Application struct {
    repository        domain.IRepository
    dispatcher        *domain.Dispatcher
    historyRepository history.IRepository
}

func NewApplication(repository *infrastructure.EventsRepository,
    historyRepository history.IRepository) (*Application, error) {
    dispatcher := domain.NewDispatcher([]domain.Subscriber{
        history.NewSubscriber(historyRepository),
    })

    return &Application{
        repository:        infrastructure.NewRepository(repository),
        historyRepository: historyRepository,
        dispatcher:        dispatcher,
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

func (a *Application) LabelTemplateHistoryList(ctx context.Context, id string) (string, error) {
    domainAggregateID, err := domain.NewLabelTemplateID(id)
    if err != nil {
        return "", err
    }

    domainHistoryList, err := a.historyRepository.List(ctx, domainAggregateID)
    if err != nil {
        return "", err
    }

    result := []LabelTemplateHistoryRow{}

    for _, domainHistory := range domainHistoryList {
        historyRow := LabelTemplateHistoryRow{
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

    resultString, err := json.Marshal(result)
    if err != nil {
        return "", err
    }

    return string(resultString), nil
}

func (a *Application) LabelTemplateAddCategoryList(ctx context.Context, labelTemplateID string,
    categoryList []string) error {
    domainCategoryList := make([]domain.Category, 0, len(categoryList))
    for _, category := range categoryList {
        categoryParts := strings.Split(category, "-")

        var categoryTypeID *int64

        if len(categoryParts) == 0 {
            return ErrCategoryEmpty
        }

        categoryID, err := strconv.ParseInt(categoryParts[0], 10, 64)
        if err != nil {
            return err
        }

        const partsLimit = 2
        if len(categoryParts) == partsLimit {
            categoryTypeIDValue, err := strconv.ParseInt(categoryParts[1], 10, 64)
            if err != nil {
                return err
            }

            categoryTypeID = &categoryTypeIDValue
        }

        domainCategory, err := domain.NewCategory(categoryID, categoryTypeID)
        if err != nil {
            return err
        }

        domainCategoryList = append(domainCategoryList, domainCategory)
    }

    domainLabel, err := a.loadLabelTemplate(ctx, labelTemplateID)
    if err != nil {
        return err
    }

    err = domainLabel.AddCategoryList(domainCategoryList)
    if err != nil {
        return err
    }

    err = a.repository.Save(ctx, domainLabel)
    if err != nil {
        return err
    }

    return nil
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
