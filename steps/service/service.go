package service

import (
    "context"
    "errors"
    "net/mail"
    "net/url"
    "strings"

    "effective-architecture/steps/contract"
)

var (
    _                                                 contract.IApplication = Service{}
    ErrLabelTemplateAlreadyCreated                                          = errors.New("попытка создать уже существующий шаблон")
    ErrLabelTemplateAlreadyDeleted                                          = errors.New("попытка удалить уже удалённый шаблон")
    ErrLabelTemplateWrongManufacturerOrganizationName                       = errors.New("название организации производителя должно " +
        "быть до 255 символов в длину")
    ErrLabelTemplateManufacturerOrganizationAddressWrongLen = errors.New("адрес должен быть до 255 символов в длину")
    ErrLabelTemplateManufacturerEmailWrongLen               = errors.New("email должен быть до 255 символов в длину")
    ErrLabelTemplateManufacturerSiteWrongLen                = errors.New("сайт должен быть до 255 символов в длину")
    ErrLabelTemplateManufacturerSiteWrongFormat             = errors.New("сайт имеет не корректный формат")
    ErrLabelTemplateManufacturerEmailWrongFormat            = errors.New("email имеет не корректный формат")
)

type IRepository interface {
    Insert(ctx context.Context, model LabelTemplate) error
    Find(ctx context.Context, id string) (LabelTemplate, error)
    Update(ctx context.Context, model LabelTemplate) error
    Truncate(ctx context.Context) error
    Delete(ctx context.Context, id string) error
}

type IHistoryRepository interface {
    Create(ctx context.Context, model LabelTemplateHistory, orderKey int) error
    FindAll(ctx context.Context, labelTemplateID string) ([]LabelTemplateHistoryResult, error)
    Truncate(ctx context.Context) error
}

type Service struct {
    repository        IRepository
    historyRepository IHistoryRepository
}

func NewService(repository IRepository, historyRepository IHistoryRepository) *Service {
    return &Service{
        repository:        repository,
        historyRepository: historyRepository,
    }
}

func (s Service) Create(ctx context.Context, labelTemplateID string,
    manufacturer contract.Manufacturer) error {
    model := LabelTemplate{
        ID:                              labelTemplateID,
        ManufacturerOrganizationName:    manufacturer.OrganizationName,
        ManufacturerOrganizationAddress: manufacturer.OrganizationAddress,
        ManufacturerEmail:               manufacturer.Email,
        ManufacturerSite:                manufacturer.Site,
    }

    err := s.validateManufacturer(manufacturer)
    if err != nil {
        return err
    }

    err = s.repository.Insert(ctx, model)
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"label_templates_pkey\"") {
            return ErrLabelTemplateAlreadyCreated
        }

        return err
    }

    action := "created"

    err = s.createHistory(ctx, labelTemplateID, manufacturer, action)
    if err != nil {
        return err
    }

    return nil
}

func (s Service) Update(ctx context.Context, labelTemplateID string,
    manufacturer contract.Manufacturer) error {
    model := LabelTemplate{
        ID:                              labelTemplateID,
        ManufacturerOrganizationName:    manufacturer.OrganizationName,
        ManufacturerOrganizationAddress: manufacturer.OrganizationAddress,
        ManufacturerEmail:               manufacturer.Email,
        ManufacturerSite:                manufacturer.Site,
    }

    err := s.validateManufacturer(manufacturer)
    if err != nil {
        return err
    }

    err = s.repository.Update(ctx, model)
    if err != nil {
        return err
    }

    action := "updated"

    err = s.createHistory(ctx, labelTemplateID, manufacturer, action)
    if err != nil {
        return err
    }

    return nil
}

func (s Service) Get(ctx context.Context, labelTemplateID string) (contract.LabelTemplate, error) {
    model, err := s.repository.Find(ctx, labelTemplateID)
    if err != nil {
        return contract.LabelTemplate{}, err
    }

    response := contract.LabelTemplate{
        ID: labelTemplateID,
        Manufacturer: contract.Manufacturer{
            OrganizationName:    model.ManufacturerOrganizationName,
            OrganizationAddress: model.ManufacturerOrganizationAddress,
            Email:               model.ManufacturerEmail,
            Site:                model.ManufacturerSite,
        },
    }

    return response, nil
}

func (s Service) Delete(ctx context.Context, labelTemplateID string) error {
    err := s.repository.Delete(ctx, labelTemplateID)
    if err != nil {
        if errors.Is(err, ErrCouldNotDelete) {
            return ErrLabelTemplateAlreadyDeleted
        }

        return err
    }

    err = s.historyRepository.Create(ctx, LabelTemplateHistory{
        LabelTemplateID: labelTemplateID,
        Action:          "deleted",
    }, 0)
    if err != nil {
        return err
    }

    return nil
}

func (s Service) HistoryList(ctx context.Context, labelTemplateID string) ([]contract.LabelTemplateHistoryRow, error) {
    historyList, err := s.historyRepository.FindAll(ctx, labelTemplateID)
    if err != nil {
        return nil, err
    }

    result := make([]contract.LabelTemplateHistoryRow, 0, len(historyList))
    for _, history := range historyList {
        result = append(result, contract.LabelTemplateHistoryRow{
            OrderKey:                           history.OrderKey,
            Action:                             history.Action,
            NewManufacturerOrganizationName:    history.NewManufacturerOrganizationName,
            NewManufacturerOrganizationAddress: history.NewManufacturerOrganizationAddress,
            NewManufacturerEmail:               history.NewManufacturerEmail,
            NewManufacturerSite:                history.NewManufacturerSite,
        })
    }

    return result, nil
}

func (s Service) AddCategoryList(ctx context.Context, labelTemplateID string, categoryList []contract.Category) error {
    return nil
}

func (s Service) validateManufacturer(manufacturer contract.Manufacturer) error {
    const varcharLimit = 255
    if len(manufacturer.OrganizationName) > varcharLimit || len(manufacturer.OrganizationName) == 0 {
        return ErrLabelTemplateWrongManufacturerOrganizationName
    }

    if len(manufacturer.OrganizationAddress) > varcharLimit {
        return ErrLabelTemplateManufacturerOrganizationAddressWrongLen
    }

    if len(manufacturer.Email) > varcharLimit {
        return ErrLabelTemplateManufacturerEmailWrongLen
    }

    if len(manufacturer.Site) > varcharLimit {
        return ErrLabelTemplateManufacturerSiteWrongLen
    }

    if len(manufacturer.Email) > 0 {
        _, err := mail.ParseAddress(manufacturer.Email)
        if err != nil {
            return ErrLabelTemplateManufacturerEmailWrongFormat
        }
    }

    if len(manufacturer.Site) > 0 {
        _, err := url.ParseRequestURI(manufacturer.Site)
        if err != nil {
            return ErrLabelTemplateManufacturerSiteWrongFormat
        }
    }

    return nil
}

func (s Service) createHistory(ctx context.Context, labelTemplateID string, manufacturer contract.Manufacturer,
    action string) error {
    err := s.historyRepository.Create(ctx, LabelTemplateHistory{
        LabelTemplateID:                    labelTemplateID,
        Action:                             action,
        NewManufacturerOrganizationName:    manufacturer.OrganizationName,
        NewManufacturerOrganizationAddress: manufacturer.OrganizationAddress,
        NewManufacturerEmail:               manufacturer.Email,
        NewManufacturerSite:                manufacturer.Site,
    }, 0)
    if err != nil {
        return err
    }

    return nil
}

type Manufacturer struct {
    OrganizationName    string `json:"organizationName"`
    OrganizationAddress string `json:"organizationAddress,omitempty"`
    Email               string `json:"email,omitempty"`
    Site                string `json:"site,omitempty"`
}

type GetLabelTemplateResponse struct {
    ID           string       `json:"id"`
    Manufacturer Manufacturer `json:"manufacturer"`
}
