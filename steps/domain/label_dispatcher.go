package domain

import "context"

type LabelSubscriber interface {
	Handle(ctx context.Context, aggregate *Label, event LabelEvent) error
}

type LabelDispatcher struct {
	subscriberList []LabelSubscriber
}

func NewLabelDispatcher(subscriberList []LabelSubscriber) *LabelDispatcher {
	return &LabelDispatcher{subscriberList: subscriberList}
}

func (d *LabelDispatcher) Dispatch(ctx context.Context, aggregate *Label) error {
	for _, subscriber := range d.subscriberList {
		for _, event := range aggregate.Events {
			err := subscriber.Handle(ctx, aggregate, event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
