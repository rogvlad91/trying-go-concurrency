package process

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

func (i *Implementation) Process(order model.Order) (model.Order, error) {
	start := time.Now().UTC()

	time.Sleep(1 * time.Second)

	order.Tracking = append(order.Tracking, model.OrderTracking{
		State: model.OrderStateProcessing,
		Start: start,
	})

	order.Storage = model.StorageId(order.Id % 2)

	return order, nil
}

func (i *Implementation) Pipeline(ctx context.Context, orderCh <-chan model.PipelineOrder) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)

	go func() {
		defer close(outCh)

		for order := range orderCh {
			orderProcessed, err := i.Process(order.Order)
			select {
			case <-ctx.Done():
				return
			case outCh <- model.PipelineOrder{
				Order: orderProcessed,
				Err:   err,
			}:
			}
		}
	}()

	return outCh
}
