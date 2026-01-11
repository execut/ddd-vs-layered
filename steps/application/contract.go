package application

type LabelTemplateHistoryList struct {
    List []LabelTemplateHistoryRow `json:"list"`
}

type LabelTemplateHistoryRow struct {
    OrderKey                           int    `json:"orderKey"`
    Action                             string `json:"action"`
    NewManufacturerOrganizationName    string `json:"newManufacturerOrganizationName,omitempty"`
    NewManufacturerOrganizationAddress string `json:"newManufacturerOrganizationAddress,omitempty"`
    NewManufacturerEmail               string `json:"newManufacturerEmail,omitempty"`
    NewManufacturerSite                string `json:"newManufacturerSite,omitempty"`
}
