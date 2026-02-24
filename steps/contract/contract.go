package contract

import (
	"context"
	"errors"
)

var ErrLabelTemplateWrongUser = errors.New("нельзя редактировать чужие записи")
var ErrLabelWrongUser = errors.New("нельзя редактировать чужие записи")

type IApplication interface { //nolint:interfacebloat
	Create(ctx context.Context, userID string, labelTemplateID string, manufacturer Manufacturer) error
	Get(ctx context.Context, userID string, labelTemplateID string) (LabelTemplate, error)
	Delete(ctx context.Context, userID string, labelTemplateID string) error
	Update(ctx context.Context, userID string, labelTemplateID string, manufacturer Manufacturer) error
	HistoryList(ctx context.Context, userID string, labelTemplateID string) ([]LabelTemplateHistoryRow, error)
	AddCategoryList(ctx context.Context, userID string, labelTemplateID string, categoryList []Category) error
	UnlinkCategoryList(ctx context.Context, userID string, labelTemplateID string, categoryList []Category) error
	Deactivate(ctx context.Context, userID string, labelTemplateID string) error
	Activate(ctx context.Context, userID string, labelTemplateID string) error
	Cleanup(ctx context.Context, userID string, labelTemplateID string) error

	StartLabelGeneration(ctx context.Context, userID string, generationID string, sku int64) error
	LabelGeneration(ctx context.Context, userID string, generationID string) (LabelGeneration, error)
	FillLabelGeneration(ctx context.Context, userID string, generationID string) error
	GenerateLabel(ctx context.Context, userID string, generationID string) error
}

type LabelTemplateHistoryList struct {
	List []LabelTemplateHistoryRow `json:"list"`
}

type LabelTemplate struct {
	ID           string       `json:"id"`
	Manufacturer Manufacturer `json:"manufacturer"`
}

type LabelTemplateHistoryRow struct {
	OrderKey                           int                           `json:"orderKey"`
	Action                             LabelTemplateHistoryRowAction `json:"action"`
	NewManufacturerOrganizationName    string                        `json:"newManufacturerOrganizationName,omitempty"`
	NewManufacturerOrganizationAddress string                        `json:"newManufacturerOrganizationAddress,omitempty"`
	NewManufacturerEmail               string                        `json:"newManufacturerEmail,omitempty"`
	NewManufacturerSite                string                        `json:"newManufacturerSite,omitempty"`
	CategoryList                       []Category                    `json:"categoryList,omitempty"`
}

type LabelTemplateHistoryRowAction string

const (
	LabelTemplateHistoryRowActionCreated              LabelTemplateHistoryRowAction = "created"
	LabelTemplateHistoryRowActionUpdated              LabelTemplateHistoryRowAction = "updated"
	LabelTemplateHistoryRowActionDeleted              LabelTemplateHistoryRowAction = "deleted"
	LabelTemplateHistoryRowActionCategoryListAdded    LabelTemplateHistoryRowAction = "category_list_added"
	LabelTemplateHistoryRowActionCategoryListUnlinked LabelTemplateHistoryRowAction = "category_list_unlinked"
	LabelTemplateHistoryRowActionActivated            LabelTemplateHistoryRowAction = "activated"
	LabelTemplateHistoryRowActionDeactivated          LabelTemplateHistoryRowAction = "deactivated"
)

type Manufacturer struct {
	OrganizationName    string `json:"organizationName"`
	OrganizationAddress string `json:"organizationAddress,omitempty"`
	Site                string `json:"site,omitempty"`
	Email               string `json:"email,omitempty"`
}

type Category struct {
	CategoryID string
	TypeID     *string
}

type CategoryWithType struct {
	CategoryID string
	TypeID     string
}

type LabelGeneration struct {
	Status   LabelGenerationStatus
	FilePath *string
}

type LabelGenerationStatus string

const (
	LabelGenerationStatusGeneration LabelGenerationStatus = "Генерация"
	LabelGenerationStatusDataFilled LabelGenerationStatus = "Наполнено данными"
	LabelGenerationStatusGenerated  LabelGenerationStatus = "Сгенерировано"
)

type Product struct {
	Name         string
	Manufacturer Manufacturer
	SKU          int64
}
