package complete

import (
	"context"
	"time"
	"trying-concurrency-go/internal/model"
)

func New() *Implementation {
	return &Implementation{}
}

type Implementation struct {
}

func (i *Implementation) Complete(order model.Order) (model.Order, error) {
	start := time.Now().UTC()

	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateCompleted,
		Start: start,
	})

	order.PickupPoint = model.PickupPointId(order.Storage) + model.PickupPointId(order.Id)

	return order, nil
}

func (i *Implementation) Pipeline(ctx context.Context, orderCh <-chan model.PipelineOrder) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)

	go func() {
		defer close(outCh)

		for order := range orderCh {
			orderCompleted, err := i.Complete(order.Order)
			select {
			case <-ctx.Done():
				return
			case outCh <- model.PipelineOrder{
				Order: orderCompleted,
				Err:   err,
			}:
			}
		}
	}()
	return outCh
}
