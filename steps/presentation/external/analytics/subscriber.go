package analytics

import (
	"context"
	"effective-architecture/steps/contract/external"

	"effective-architecture/steps/domain"
)

type Subscriber struct {
	analytics external.IAnalytics
}

func NewSubscriber(analytics external.IAnalytics) *Subscriber {
	return &Subscriber{analytics: analytics}
}

func (s Subscriber) Handle(ctx context.Context,
	aggregate *domain.LabelTemplate, event domain.LabelTemplateEvent) error {
	var (
		userID    = aggregate.UserID
		eventType external.AnalyticsEventType
	)

	switch event.(type) {
	case domain.LabelTemplateCreatedEvent:
		eventType = external.AnalyticsEventTypeCreated
	case domain.LabelTemplateUpdatedEvent:
		eventType = external.AnalyticsEventTypeUpdated
	case domain.LabelTemplateDeletedEvent:
		eventType = external.AnalyticsEventTypeDeleted
	case domain.LabelTemplateCategoryListAddedEvent:
		eventType = external.AnalyticsEventTypeCategoryListAdded
	case domain.LabelTemplateCategoryListUnlinkedEvent:
		eventType = external.AnalyticsEventTypeCategoryListUnlinked
	case domain.LabelTemplateActivatedEvent:
		eventType = external.AnalyticsEventTypeActivated
	case domain.LabelTemplateDeactivatedEvent:
		eventType = external.AnalyticsEventTypeDeactivated
	}

	if eventType == "" {
		return nil
	}

	err := s.analytics.Send(ctx, eventType, userID)
	if err != nil {
		return err
	}

	return nil
}
