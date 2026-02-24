package analytics

import (
	"context"
	"effective-architecture/steps/contract/external"

	"effective-architecture/steps/domain"
)

type LabelSubscriber struct {
	analytics external.IAnalytics
}

func NewLabelSubscriber(analytics external.IAnalytics) *LabelSubscriber {
	return &LabelSubscriber{analytics: analytics}
}

func (s LabelSubscriber) Handle(ctx context.Context,
	aggregate *domain.Label, event domain.LabelEvent) error {
	var (
		userID    = aggregate.UserID
		eventType external.AnalyticsEventType
	)

	switch event.(type) {
	case domain.LabelGenerationStartedEvent:
		eventType = external.AnalyticsEventTypeLabelGenerationStarted
	case domain.LabelDataFilledEvent:
		eventType = external.AnalyticsEventTypeLabelGenerationFilled
	case domain.LabelGeneratedEvent:
		eventType = external.AnalyticsEventTypeLabelGenerated
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
