package history

import (
	"effective-architecture/steps/contract"
	"errors"

	"effective-architecture/steps/domain"

	"github.com/google/uuid"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
)

type History struct {
	ID                                 uuid.UUID
	AggregateID                        domain.LabelTemplateID
	OrderKey                           int
	Action                             contract.LabelTemplateHistoryRowAction
	NewManufacturerOrganizationName    *domain.OrganizationName
	NewManufacturerOrganizationAddress *domain.OrganizationAddress
	NewManufacturerEmail               *domain.Email
	NewManufacturerSite                *domain.Site
	CategoryList                       []domain.Category
}

func NewHistory(aggregateID domain.LabelTemplateID, orderKey int, action contract.LabelTemplateHistoryRowAction,
	newManufacturerOrganizationName *domain.OrganizationName,
	newManufacturerOrganizationAddress *domain.OrganizationAddress, newManufacturerEmail *domain.Email,
	newManufacturerSite *domain.Site, categoryList []domain.Category) (History, error) {
	return History{AggregateID: aggregateID, OrderKey: orderKey, Action: action,
		NewManufacturerOrganizationName:    newManufacturerOrganizationName,
		NewManufacturerOrganizationAddress: newManufacturerOrganizationAddress,
		NewManufacturerEmail:               newManufacturerEmail,
		NewManufacturerSite:                newManufacturerSite,
		CategoryList:                       categoryList,
	}, nil
}

func NewHistoryFromEvent(aggregate *domain.LabelTemplate,
	event domain.LabelTemplateEvent, currentCount int) (History, error) {
	var (
		action                             contract.LabelTemplateHistoryRowAction
		newManufacturerOrganizationName    *domain.OrganizationName
		newManufacturerOrganizationAddress *domain.OrganizationAddress
		newManufacturerEmail               *domain.Email
		newManufacturerSite                *domain.Site
		categoryList                       []domain.Category
	)

	switch payload := event.(type) {
	case domain.LabelTemplateCreatedEvent:
		action = contract.LabelTemplateHistoryRowActionCreated
		newManufacturerOrganizationName = &payload.Manufacturer.OrganizationName
		newManufacturerOrganizationAddress = payload.Manufacturer.OrganizationAddress
		newManufacturerEmail = payload.Manufacturer.Email
		newManufacturerSite = payload.Manufacturer.Site
	case domain.LabelTemplateDeletedEvent:
		action = contract.LabelTemplateHistoryRowActionDeleted
	case domain.LabelTemplateUpdatedEvent:
		action = contract.LabelTemplateHistoryRowActionUpdated
		newManufacturerOrganizationName = &payload.Manufacturer.OrganizationName
		newManufacturerOrganizationAddress = payload.Manufacturer.OrganizationAddress
		newManufacturerEmail = payload.Manufacturer.Email
		newManufacturerSite = payload.Manufacturer.Site
	case domain.LabelTemplateCategoryListAddedEvent:
		action = contract.LabelTemplateHistoryRowActionCategoryListAdded
		categoryList = payload.CategoryList
	case domain.LabelTemplateCategoryListUnlinkedEvent:
		action = contract.LabelTemplateHistoryRowActionCategoryListUnlinked
		categoryList = payload.CategoryList
	case domain.LabelTemplateActivatedEvent:
		action = contract.LabelTemplateHistoryRowActionActivated
	case domain.LabelTemplateDeactivatedEvent:
		action = contract.LabelTemplateHistoryRowActionDeactivated
	default:
		return History{}, ErrUnknownEventType
	}

	return History{
		ID:                                 uuid.New(),
		AggregateID:                        aggregate.ID,
		OrderKey:                           currentCount + 1,
		Action:                             action,
		NewManufacturerOrganizationName:    newManufacturerOrganizationName,
		NewManufacturerOrganizationAddress: newManufacturerOrganizationAddress,
		NewManufacturerEmail:               newManufacturerEmail,
		NewManufacturerSite:                newManufacturerSite,
		CategoryList:                       categoryList,
	}, nil
}
