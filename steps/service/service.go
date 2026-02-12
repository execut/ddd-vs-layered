package service

import (
    "context"
    "errors"
    "net/mail"
    "net/url"
    "strconv"
    "strings"

    "effective-architecture/steps/contract"
    "effective-architecture/steps/contract/external"
)

var (
    _ contract.IApplication = (*Service)(nil)

    ErrLabelTemplateAlreadyCreated                    = errors.New("попытка создать уже существующий шаблон")
    ErrLabelTemplateAlreadyDeleted                    = errors.New("попытка удалить уже удалённый шаблон")
    ErrLabelTemplateWrongManufacturerOrganizationName = errors.New("название организации производителя должно " +
        "быть до 255 символов в длину")
    ErrLabelTemplateManufacturerOrganizationAddressWrongLen = errors.New("адрес должен быть до 255 символов в длину")
    ErrLabelTemplateManufacturerEmailWrongLen               = errors.New("email должен быть до 255 символов в длину")
    ErrLabelTemplateManufacturerSiteWrongLen                = errors.New("сайт должен быть до 255 символов в длину")
    ErrLabelTemplateManufacturerSiteWrongFormat             = errors.New("сайт имеет не корректный формат")
    ErrLabelTemplateManufacturerEmailWrongFormat            = errors.New("email имеет не корректный формат")
)

type Service struct {
    repository                        *Repository
    historyRepository                 *HistoryRepository
    vsCategoryRepository              *VsCategoryRepository
    ozonService                       external.IExternalServiceOzon
    categoryVsLabelTemplateRepository *CategoryVsLabelTemplateRepository
    labelRepository                   *LabelRepository
}

func NewService(repository *Repository, historyRepository *HistoryRepository,
    categoryRepository *VsCategoryRepository, ozonService external.IExternalServiceOzon,
    categoryVsLabelTemplateRepository *CategoryVsLabelTemplateRepository, labelRepository *LabelRepository) *Service {
    return &Service{
        repository:                        repository,
        historyRepository:                 historyRepository,
        vsCategoryRepository:              categoryRepository,
        ozonService:                       ozonService,
        categoryVsLabelTemplateRepository: categoryVsLabelTemplateRepository,
        labelRepository:                   labelRepository,
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
        orderKey := history.OrderKey

        orderKeyAsInt, err := strconv.Atoi(orderKey)
        if err != nil {
            return nil, err
        }

        row := contract.LabelTemplateHistoryRow{
            OrderKey:                           orderKeyAsInt,
            Action:                             history.Action,
            NewManufacturerOrganizationName:    history.NewManufacturerOrganizationName,
            NewManufacturerOrganizationAddress: history.NewManufacturerOrganizationAddress,
            NewManufacturerEmail:               history.NewManufacturerEmail,
            NewManufacturerSite:                history.NewManufacturerSite,
        }

        if len(history.CategoryList) > 0 {
            row.CategoryList = mapHistoryCategoryToContract(history.CategoryList)
        }

        result = append(result, row)
    }

    return result, nil
}

func (s Service) AddCategoryList(ctx context.Context, labelTemplateID string, categoryList []contract.Category) error {
    vsCategoryModelList, err := createVsCategoryModelList(categoryList, labelTemplateID)
    if err != nil {
        return err
    }

    for _, vsCategoryModel := range vsCategoryModelList {
        err := s.vsCategoryRepository.Create(ctx, vsCategoryModel)
        if err != nil {
            return err
        }

        model := CategoryIDVsLabelTemplateID{
            LabelTemplateID: labelTemplateID,
            CategoryID:      vsCategoryModel.CategoryID,
            TypeID:          vsCategoryModel.TypeID,
        }

        err = s.categoryVsLabelTemplateRepository.Create(ctx, model)
        if err != nil {
            return err
        }
    }

    serviceCategoryList, err := mapHistoryCategoryToService(categoryList)
    if err != nil {
        return err
    }

    err = s.historyRepository.Create(ctx, LabelTemplateHistory{
        LabelTemplateID: labelTemplateID,
        Action:          "category_list_added",
        CategoryList:    serviceCategoryList,
    }, 0)
    if err != nil {
        return err
    }

    return nil
}

func (s Service) UnlinkCategoryList(ctx context.Context, labelTemplateID string,
    categoryList []contract.Category) error {
    vsCategoryModelList, err := createVsCategoryModelList(categoryList, labelTemplateID)
    if err != nil {
        return err
    }

    for _, vsCategoryModel := range vsCategoryModelList {
        err := s.vsCategoryRepository.Delete(ctx, vsCategoryModel)
        if err != nil {
            return err
        }

        model := CategoryIDVsLabelTemplateID{
            LabelTemplateID: labelTemplateID,
            CategoryID:      vsCategoryModel.CategoryID,
            TypeID:          vsCategoryModel.TypeID,
        }

        err = s.categoryVsLabelTemplateRepository.Delete(ctx, model)
        if err != nil {
            return err
        }
    }

    serviceCategoryList, err := mapHistoryCategoryToService(categoryList)
    if err != nil {
        return err
    }

    err = s.historyRepository.Create(ctx, LabelTemplateHistory{
        LabelTemplateID: labelTemplateID,
        Action:          "category_list_unlinked",
        CategoryList:    serviceCategoryList,
    }, 0)
    if err != nil {
        return err
    }

    return nil
}

func (s Service) StartLabelGeneration(ctx context.Context, labelID string, sku int64) error {
    err := s.labelRepository.Exists(ctx, labelID)
    if err != nil {
        return err
    }

    product, err := s.ozonService.Product(ctx, sku)
    if err != nil {
        return err
    }

    templateID, err := s.categoryVsLabelTemplateRepository.LabelTemplateID(ctx, product)
    if err != nil {
        return err
    }

    label := Label{
        ID:              labelID,
        LabelTemplateID: templateID,
        SKU:             sku,
    }

    err = s.labelRepository.Create(ctx, label)
    if err != nil {
        return err
    }

    return nil
}

func (s Service) Cleanup(ctx context.Context, labelTemplateID string) error {
    _ = s.repository.Delete(ctx, labelTemplateID)
    _ = s.historyRepository.Delete(ctx, labelTemplateID)
    _ = s.vsCategoryRepository.DeleteByLabelTemplateID(ctx, labelTemplateID)
    _ = s.labelRepository.Delete(ctx, labelTemplateID)
    _ = s.categoryVsLabelTemplateRepository.DeleteByLabelTemplateID(ctx, labelTemplateID)

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

func createVsCategoryModelList(categoryList []contract.Category,
    labelTemplateID string) ([]LabelTemplateVsCategory, error) {
    vsCategoryModelList := make([]LabelTemplateVsCategory, 0, len(categoryList))

    for _, category := range categoryList {
        categoryID, err := strconv.ParseInt(category.CategoryID, 10, 64)
        if err != nil {
            return nil, err
        }

        var typeID *int64

        if category.TypeID != nil {
            typeIDValue, err := strconv.ParseInt(*category.TypeID, 10, 64)
            if err != nil {
                return nil, err
            }

            typeID = &typeIDValue
        }

        vsCategoryModel := LabelTemplateVsCategory{
            LabelTemplateID: labelTemplateID,
            CategoryID:      categoryID,
            TypeID:          typeID,
        }

        vsCategoryModelList = append(vsCategoryModelList, vsCategoryModel)
    }

    return vsCategoryModelList, nil
}

func mapHistoryCategoryToService(categoryList []contract.Category) ([]HistoryCategory, error) {
    serviceCategoryList := make([]HistoryCategory, 0, len(categoryList))
    for _, category := range categoryList {
        categoryID, err := strconv.ParseInt(category.CategoryID, 10, 64)
        if err != nil {
            return nil, err
        }

        var typeID *int64

        if category.TypeID != nil {
            typeIDValue, err := strconv.ParseInt(*category.TypeID, 10, 64)
            if err != nil {
                return nil, err
            }

            typeID = &typeIDValue
        }

        serviceCategoryList = append(serviceCategoryList, HistoryCategory{
            CategoryID: categoryID,
            TypeID:     typeID,
        })
    }

    return serviceCategoryList, nil
}

func mapHistoryCategoryToContract(serviceCategoryList []HistoryCategory) []contract.Category {
    categoryList := make([]contract.Category, 0, len(serviceCategoryList))
    for _, serviceCategory := range serviceCategoryList {
        var typeID *string

        if serviceCategory.TypeID != nil {
            typeIDValue := strconv.FormatInt(*serviceCategory.TypeID, 10)
            typeID = &typeIDValue
        }

        categoryList = append(categoryList, contract.Category{
            CategoryID: strconv.FormatInt(serviceCategory.CategoryID, 10),
            TypeID:     typeID,
        })
    }

    return categoryList
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
