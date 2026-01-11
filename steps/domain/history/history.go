package history

import (
    "errors"

    "effective-architecture/steps/domain"
)

var (
    ErrUnknownEventType = errors.New("unknown event type")
)

type History struct {
    AggregateID                        domain.LabelTemplateID
    OrderKey                           int
    Action                             string
    NewManufacturerOrganizationName    *domain.OrganizationName
    NewManufacturerOrganizationAddress *domain.OrganizationAddress
    NewManufacturerEmail               *domain.Email
    NewManufacturerSite                *domain.Site
}

func NewHistory(aggregateID domain.LabelTemplateID, orderKey int, action string,
    newManufacturerOrganizationName *domain.OrganizationName,
    newManufacturerOrganizationAddress *domain.OrganizationAddress, newManufacturerEmail *domain.Email,
    newManufacturerSite *domain.Site) (History, error) {
    return History{AggregateID: aggregateID, OrderKey: orderKey, Action: action,
        NewManufacturerOrganizationName:    newManufacturerOrganizationName,
        NewManufacturerOrganizationAddress: newManufacturerOrganizationAddress,
        NewManufacturerEmail:               newManufacturerEmail,
        NewManufacturerSite:                newManufacturerSite,
    }, nil
}

func NewHistoryFromEvent(aggregate *domain.LabelTemplate,
    event domain.LabelTemplateEvent, currentCount int) (History, error) {
    var (
        action                             string
        newManufacturerOrganizationName    *domain.OrganizationName
        newManufacturerOrganizationAddress *domain.OrganizationAddress
        newManufacturerEmail               *domain.Email
        newManufacturerSite                *domain.Site
    )

    switch payload := event.(type) {
    case domain.LabelTemplateCreatedEvent:
        action = "created"
        newManufacturerOrganizationName = &payload.Manufacturer.OrganizationName
        newManufacturerOrganizationAddress = payload.Manufacturer.OrganizationAddress
        newManufacturerEmail = payload.Manufacturer.Email
        newManufacturerSite = payload.Manufacturer.Site
    case domain.LabelTemplateDeletedEvent:
        action = "deleted"
    case domain.LabelTemplateUpdatedEvent:
        action = "updated"
        newManufacturerOrganizationName = &payload.Manufacturer.OrganizationName
        newManufacturerOrganizationAddress = payload.Manufacturer.OrganizationAddress
        newManufacturerEmail = payload.Manufacturer.Email
        newManufacturerSite = payload.Manufacturer.Site
    default:
        return History{}, ErrUnknownEventType
    }

    return History{
        AggregateID:                        aggregate.ID,
        OrderKey:                           currentCount + 1,
        Action:                             action,
        NewManufacturerOrganizationName:    newManufacturerOrganizationName,
        NewManufacturerOrganizationAddress: newManufacturerOrganizationAddress,
        NewManufacturerEmail:               newManufacturerEmail,
        NewManufacturerSite:                newManufacturerSite,
    }, nil
}
