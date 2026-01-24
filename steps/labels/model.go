package labels

type LabelTemplate struct {
    ID                              string
    ManufacturerOrganizationName    string
    ManufacturerOrganizationAddress string
    ManufacturerEmail               string
    ManufacturerSite                string
}

type LabelTemplateHistory struct {
    LabelTemplateID                    string
    OrderKey                           string
    Action                             string
    NewManufacturerOrganizationName    string
    NewManufacturerOrganizationAddress string
    NewManufacturerEmail               string
    NewManufacturerSite                string
}

type LabelTemplateHistoryResult struct {
    OrderKey                           int    `json:"orderKey"`
    Action                             string `json:"action"`
    NewManufacturerOrganizationName    string `json:"newManufacturerOrganizationName,omitempty"`
    NewManufacturerOrganizationAddress string `json:"newManufacturerOrganizationAddress,omitempty"`
    NewManufacturerEmail               string `json:"newManufacturerEmail,omitempty"`
    NewManufacturerSite                string `json:"newManufacturerSite,omitempty"`
}
