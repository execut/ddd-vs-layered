//go:generate mockgen -package external -destination=analytics_mocks.go . IAnalytics
package external

import "context"

type AnalyticsEventType string

const (
	AnalyticsEventTypeCreated              AnalyticsEventType = "Created"
	AnalyticsEventTypeDeleted              AnalyticsEventType = "Deleted"
	AnalyticsEventTypeUpdated              AnalyticsEventType = "Updated"
	AnalyticsEventTypeCategoryListAdded    AnalyticsEventType = "CategoryListAdded"
	AnalyticsEventTypeCategoryListUnlinked AnalyticsEventType = "CategoryListUnlinked"
	AnalyticsEventTypeActivated            AnalyticsEventType = "Activated"
	AnalyticsEventTypeDeactivated          AnalyticsEventType = "Deactivated"

	AnalyticsEventTypeLabelGenerationStarted AnalyticsEventType = "LabelGenerationStarted"
	AnalyticsEventTypeLabelGenerated         AnalyticsEventType = "LabelGenerated"
	AnalyticsEventTypeLabelGenerationFilled  AnalyticsEventType = "LabelGenerationFilled"
)

type IAnalytics interface {
	Send(ctx context.Context, event AnalyticsEventType, userID string)
}
