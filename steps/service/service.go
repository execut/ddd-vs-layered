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

	"github.com/jackc/pgx/v5"
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
	generatorService                  external.ILabelGenerator
}

func NewService(repository *Repository, historyRepository *HistoryRepository,
	categoryRepository *VsCategoryRepository, ozonService external.IExternalServiceOzon,
	categoryVsLabelTemplateRepository *CategoryVsLabelTemplateRepository, labelRepository *LabelRepository,
	generatorService external.ILabelGenerator) *Service {
	return &Service{
		repository:                        repository,
		historyRepository:                 historyRepository,
		vsCategoryRepository:              categoryRepository,
		ozonService:                       ozonService,
		categoryVsLabelTemplateRepository: categoryVsLabelTemplateRepository,
		labelRepository:                   labelRepository,
		generatorService:                  generatorService,
	}
}

func (s Service) Create(ctx context.Context, userID, labelTemplateID string,
	manufacturer contract.Manufacturer) error {
	err := s.validateManufacturer(manufacturer)
	if err != nil {
		return err
	}

	hasOldModel := false

	oldModel, err := s.repository.Find(ctx, labelTemplateID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	} else {
		hasOldModel = true

		if oldModel.Status != LabelTemplateStatusDeleted {
			return ErrLabelTemplateAlreadyCreated
		}
	}

	model := LabelTemplate{
		ID:                              labelTemplateID,
		Status:                          LabelTemplateStatusCreated,
		UserID:                          userID,
		ManufacturerOrganizationName:    manufacturer.OrganizationName,
		ManufacturerOrganizationAddress: manufacturer.OrganizationAddress,
		ManufacturerEmail:               manufacturer.Email,
		ManufacturerSite:                manufacturer.Site,
	}

	if hasOldModel {
		err = s.repository.Update(ctx, model)
	} else {
		err = s.repository.Insert(ctx, model)
	}

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"label_templates_pkey\"") {
			return ErrLabelTemplateAlreadyCreated
		}

		return err
	}

	err = s.createHistory(ctx, labelTemplateID, manufacturer, contract.LabelTemplateHistoryRowActionCreated)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Update(ctx context.Context, userID, labelTemplateID string,
	manufacturer contract.Manufacturer) error {
	model, err := s.repository.Find(ctx, labelTemplateID)
	if err != nil {
		return err
	}

	err = checkLabelTemplateUser(userID, model)
	if err != nil {
		return err
	}

	model = LabelTemplate{
		ID:                              labelTemplateID,
		UserID:                          userID,
		ManufacturerOrganizationName:    manufacturer.OrganizationName,
		ManufacturerOrganizationAddress: manufacturer.OrganizationAddress,
		ManufacturerEmail:               manufacturer.Email,
		ManufacturerSite:                manufacturer.Site,
	}

	err = s.validateManufacturer(manufacturer)
	if err != nil {
		return err
	}

	err = s.repository.Update(ctx, model)
	if err != nil {
		return err
	}

	action := contract.LabelTemplateHistoryRowActionUpdated

	err = s.createHistory(ctx, labelTemplateID, manufacturer, action)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Get(ctx context.Context, userID, labelTemplateID string) (contract.LabelTemplate, error) {
	model, err := s.repository.Find(ctx, labelTemplateID)
	if err != nil {
		return contract.LabelTemplate{}, err
	}

	err = checkLabelTemplateUser(userID, model)
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

func (s Service) Delete(ctx context.Context, userID, labelTemplateID string) error {
	model, err := s.repository.Find(ctx, labelTemplateID)
	if err != nil {
		return err
	}

	if model.Status == LabelTemplateStatusDeleted {
		return ErrLabelTemplateAlreadyDeleted
	}

	err = checkLabelTemplateUser(userID, model)
	if err != nil {
		return err
	}

	model.Status = LabelTemplateStatusDeleted

	err = s.repository.Update(ctx, model)
	if err != nil {
		return err
	}

	err = s.historyRepository.Create(ctx, LabelTemplateHistory{
		LabelTemplateID: labelTemplateID,
		Action:          contract.LabelTemplateHistoryRowActionDeleted,
	}, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) HistoryList(ctx context.Context, userID,
	labelTemplateID string) ([]contract.LabelTemplateHistoryRow, error) {
	model, err := s.repository.Find(ctx, labelTemplateID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	err = checkLabelTemplateUser(userID, model)
	if err != nil {
		return nil, err
	}

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

func (s Service) AddCategoryList(ctx context.Context, userID, labelTemplateID string,
	categoryList []contract.Category) error {
	action := contract.LabelTemplateHistoryRowActionCategoryListAdded

	return s.editCategoryList(ctx, userID, labelTemplateID, categoryList, true, action)
}

func (s Service) UnlinkCategoryList(ctx context.Context, userID, labelTemplateID string,
	categoryList []contract.Category) error {
	action := contract.LabelTemplateHistoryRowActionCategoryListUnlinked

	return s.editCategoryList(ctx, userID, labelTemplateID, categoryList, false, action)
}

func (s Service) Deactivate(ctx context.Context, userID, labelTemplateID string) error {
	return s.changeLabelTemplateActivity(ctx, userID, labelTemplateID, false)
}

func (s Service) Activate(ctx context.Context, userID, labelTemplateID string) error {
	return s.changeLabelTemplateActivity(ctx, userID, labelTemplateID, true)
}

func (s Service) StartLabelGeneration(ctx context.Context, userID, labelID string, sku int64) error {
	err := s.labelRepository.Exists(ctx, labelID)
	if err != nil {
		return err
	}

	label := Label{
		ID:     labelID,
		UserID: userID,
		SKU:    sku,
		Status: contract.LabelGenerationStatusGeneration,
	}

	err = s.labelRepository.Create(ctx, label)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) LabelGeneration(ctx context.Context, userID, id string) (contract.LabelGeneration, error) {
	label, err := s.labelRepository.Get(ctx, id)
	if err != nil {
		return contract.LabelGeneration{}, err
	}

	err = checkLabelUser(userID, label)
	if err != nil {
		return contract.LabelGeneration{}, err
	}

	return contract.LabelGeneration{
		Status:   label.Status,
		FilePath: label.File,
	}, nil
}

func (s Service) FillLabelGeneration(ctx context.Context, userID, generationID string) error {
	label, err := s.labelRepository.Get(ctx, generationID)
	if err != nil {
		return err
	}

	err = checkLabelUser(userID, label)
	if err != nil {
		return err
	}

	product, err := s.ozonService.Product(ctx, label.SKU)
	if err != nil {
		return err
	}

	templateID, err := s.categoryVsLabelTemplateRepository.LabelTemplateID(ctx, product)
	if err != nil {
		return err
	}

	template, err := s.repository.Find(ctx, templateID)
	if err != nil {
		return err
	}

	label.LabelTemplateID = &templateID
	label.Status = contract.LabelGenerationStatusDataFilled
	label.ProductName = &product.Name
	label.ManufacturerOrganizationName = &template.ManufacturerOrganizationName
	label.ManufacturerOrganizationAddress = &template.ManufacturerOrganizationAddress
	label.ManufacturerEmail = &template.ManufacturerEmail
	label.ManufacturerSite = &template.ManufacturerSite

	err = s.labelRepository.Update(ctx, label)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GenerateLabel(ctx context.Context, userID, generationID string) error {
	label, err := s.labelRepository.Get(ctx, generationID)
	if err != nil {
		return err
	}

	err = checkLabelUser(userID, label)
	if err != nil {
		return err
	}

	label.Status = contract.LabelGenerationStatusGenerated

	manufacturer := contract.Manufacturer{}
	if label.ManufacturerOrganizationName != nil {
		manufacturer.OrganizationName = *label.ManufacturerOrganizationName
	}

	if label.ManufacturerOrganizationAddress != nil {
		manufacturer.OrganizationAddress = *label.ManufacturerOrganizationAddress
	}

	if label.ManufacturerEmail != nil {
		manufacturer.Email = *label.ManufacturerEmail
	}

	if label.ManufacturerSite != nil {
		manufacturer.Site = *label.ManufacturerSite
	}

	contractProduct := contract.Product{
		Name:         *label.ProductName,
		Manufacturer: manufacturer,
		SKU:          label.SKU,
	}

	generatorFile, err := s.generatorService.Generate(ctx, contractProduct)
	if err != nil {
		return err
	}

	label.File = &generatorFile.Path

	err = s.labelRepository.Update(ctx, label)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Cleanup(ctx context.Context, userID, labelTemplateID string) error {
	_ = userID
	_ = s.repository.Delete(ctx, labelTemplateID)
	_ = s.historyRepository.Delete(ctx, labelTemplateID)
	_ = s.vsCategoryRepository.DeleteByLabelTemplateID(ctx, labelTemplateID)
	_ = s.labelRepository.Delete(ctx, labelTemplateID)
	_ = s.categoryVsLabelTemplateRepository.DeleteByLabelTemplateID(ctx, labelTemplateID)

	return nil
}

func (s Service) changeLabelTemplateActivity(ctx context.Context, userID, labelTemplateID string, isActive bool) error {
	model, err := s.repository.Find(ctx, labelTemplateID)
	if err != nil {
		return err
	}

	err = checkLabelTemplateUser(userID, model)
	if err != nil {
		return err
	}

	model.IsActive = isActive

	err = s.repository.Update(ctx, model)
	if err != nil {
		return err
	}

	var action contract.LabelTemplateHistoryRowAction
	if isActive {
		action = contract.LabelTemplateHistoryRowActionActivated
	} else {
		action = contract.LabelTemplateHistoryRowActionDeactivated
	}

	err = s.historyRepository.Create(ctx, LabelTemplateHistory{
		LabelTemplateID: labelTemplateID,
		Action:          action,
	}, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) editCategoryList(ctx context.Context, userID string, labelTemplateID string,
	categoryList []contract.Category, isAdd bool, action contract.LabelTemplateHistoryRowAction) error {
	model, err := s.repository.Find(ctx, labelTemplateID)
	if err != nil {
		return err
	}

	err = checkLabelTemplateUser(userID, model)
	if err != nil {
		return err
	}

	vsCategoryModelList, err := createVsCategoryModelList(categoryList, labelTemplateID)
	if err != nil {
		return err
	}

	for _, vsCategoryModel := range vsCategoryModelList {
		if isAdd {
			err = s.vsCategoryRepository.Create(ctx, vsCategoryModel)
		} else {
			err = s.vsCategoryRepository.Delete(ctx, vsCategoryModel)
		}

		if err != nil {
			return err
		}

		model := CategoryIDVsLabelTemplateID{
			LabelTemplateID: labelTemplateID,
			CategoryID:      vsCategoryModel.CategoryID,
			TypeID:          vsCategoryModel.TypeID,
		}

		if isAdd {
			err = s.categoryVsLabelTemplateRepository.Create(ctx, model)
		} else {
			err = s.categoryVsLabelTemplateRepository.Delete(ctx, model)
		}

		if err != nil {
			return err
		}
	}

	return s.createHistoryForCategoryList(ctx, categoryList, labelTemplateID, action)
}

func (s Service) createHistoryForCategoryList(ctx context.Context, categoryList []contract.Category,
	labelTemplateID string, action contract.LabelTemplateHistoryRowAction) error {
	serviceCategoryList, err := mapHistoryCategoryToService(categoryList)
	if err != nil {
		return err
	}

	err = s.historyRepository.Create(ctx, LabelTemplateHistory{
		LabelTemplateID: labelTemplateID,
		Action:          action,
		CategoryList:    serviceCategoryList,
	}, 0)
	if err != nil {
		return err
	}

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
	action contract.LabelTemplateHistoryRowAction) error {
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

func checkLabelUser(userID string, label Label) error {
	if label.UserID != userID {
		return contract.ErrLabelWrongUser
	}

	return nil
}

func checkLabelTemplateUser(userID string, label LabelTemplate) error {
	if label.UserID != userID {
		return contract.ErrLabelTemplateWrongUser
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
