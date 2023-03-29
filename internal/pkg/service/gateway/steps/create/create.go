package create

import (
	"context"
	"homework/internal/model"
	"time"
)

func New() *Implementation {
	return &Implementation{}
}

type Implementation struct {
}

func (i *Implementation) Create(workerId model.WorkerId, orderId model.OrderId) (model.Order, error) {
	start := time.Now().UTC()

	order := model.Order{
		Id:     orderId,
		Worker: workerId,
		Tracking: []model.OrderTracking{{
			State: model.OrderStateCreated,
			Start: start,
		}},
	}

	return order, nil
}

func (i *Implementation) Pipeline(ctx context.Context, workerId model.WorkerId, orderIdChan <-chan model.OrderId) <-chan model.PipelineOrder {
	outCh := make(chan model.PipelineOrder)

	go func() {
		defer close(outCh)

		for orderId := range orderIdChan {
			order, err := i.Create(workerId, orderId)
			select {
			case <-ctx.Done():
				return
			case outCh <- model.PipelineOrder{
				Order: order,
				Err:   err,
			}:
			}
		}
	}()

	return outCh
}
