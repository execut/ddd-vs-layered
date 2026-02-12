package contract

import (
    "context"
)

type IApplication interface {
    Create(ctx context.Context, labelTemplateID string, manufacturer Manufacturer) error
    Get(ctx context.Context, labelTemplateID string) (LabelTemplate, error)
    Delete(ctx context.Context, labelTemplateID string) error
    Update(ctx context.Context, labelTemplateID string, manufacturer Manufacturer) error
    HistoryList(ctx context.Context, labelTemplateID string) ([]LabelTemplateHistoryRow, error)
    AddCategoryList(ctx context.Context, labelTemplateID string, categoryList []Category) error
    Cleanup(ctx context.Context, labelTemplateID string) error

    StartLabelGeneration(ctx context.Context, labelTemplateID string, sku int64) error
}

type LabelTemplateHistoryList struct {
    List []LabelTemplateHistoryRow `json:"list"`
}

type LabelTemplate struct {
    ID           string       `json:"id"`
    Manufacturer Manufacturer `json:"manufacturer"`
}

type LabelTemplateHistoryRow struct {
    OrderKey                           int        `json:"orderKey"`
    Action                             string     `json:"action"`
    NewManufacturerOrganizationName    string     `json:"newManufacturerOrganizationName,omitempty"`
    NewManufacturerOrganizationAddress string     `json:"newManufacturerOrganizationAddress,omitempty"`
    NewManufacturerEmail               string     `json:"newManufacturerEmail,omitempty"`
    NewManufacturerSite                string     `json:"newManufacturerSite,omitempty"`
    CategoryList                       []Category `json:"categoryList,omitempty"`
}

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
