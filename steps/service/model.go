package service

import (
	"effective-architecture/steps/contract"

	"github.com/google/uuid"
)

type LabelTemplate struct {
	ID                              string
	ManufacturerOrganizationName    string
	ManufacturerOrganizationAddress string
	ManufacturerEmail               string
	ManufacturerSite                string
}

type LabelTemplateHistory struct {
	ID                                 uuid.UUID
	LabelTemplateID                    string
	OrderKey                           string
	Action                             string
	NewManufacturerOrganizationName    string
	NewManufacturerOrganizationAddress string
	NewManufacturerEmail               string
	NewManufacturerSite                string
	CategoryList                       []HistoryCategory
}

type HistoryCategory struct {
	CategoryID int64
	TypeID     *int64
}

type LabelTemplateHistoryResult struct {
	OrderKey                           int    `json:"orderKey"`
	Action                             string `json:"action"`
	NewManufacturerOrganizationName    string `json:"newManufacturerOrganizationName,omitempty"`
	NewManufacturerOrganizationAddress string `json:"newManufacturerOrganizationAddress,omitempty"`
	NewManufacturerEmail               string `json:"newManufacturerEmail,omitempty"`
	NewManufacturerSite                string `json:"newManufacturerSite,omitempty"`
	CategoryList                       []HistoryCategory
}

type LabelTemplateVsCategory struct {
	LabelTemplateID string
	CategoryID      int64
	TypeID          *int64
}

type CategoryIDVsLabelTemplateID struct {
	CategoryID      int64
	TypeID          *int64
	LabelTemplateID string
}

type Label struct {
	ID              string
	LabelTemplateID string
	SKU             int64
	Status          contract.LabelGenerationStatus
}
