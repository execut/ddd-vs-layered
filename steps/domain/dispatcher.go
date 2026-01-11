package domain

import "context"

type Subscriber interface {
    Handle(ctx context.Context, aggregate *LabelTemplate, event LabelTemplateEvent) error
}

type Dispatcher struct {
    subscriberList []Subscriber
}

func NewDispatcher(subscriberList []Subscriber) *Dispatcher {
    return &Dispatcher{subscriberList: subscriberList}
}

func (d *Dispatcher) Dispatch(ctx context.Context, aggregate *LabelTemplate) error {
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
