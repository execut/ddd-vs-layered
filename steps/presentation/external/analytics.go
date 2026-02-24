package external

import (
	"context"
	"effective-architecture/steps/contract/external"
)

type Analytics struct {
}

func NewAnalytics() Analytics {
	return Analytics{}
}

func (a Analytics) Send(_ context.Context, _ external.AnalyticsEventType, _ string) error {
	return nil
}
