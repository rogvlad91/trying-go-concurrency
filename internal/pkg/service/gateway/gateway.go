package gateway

import (
	"context"
	"homework/internal/model"
	completestep "homework/internal/pkg/service/gateway/steps/complete"
	createstep "homework/internal/pkg/service/gateway/steps/create"
	processstep "homework/internal/pkg/service/gateway/steps/process"
	"log"
	"sync"
)

func New(create *createstep.Implementation, process *processstep.Implementation, complete *completestep.Implementation) *Implementation {
	return &Implementation{
		create:   create,
		process:  process,
		complete: complete,
	}
}

type Implementation struct {
	create   *createstep.Implementation
	process  *processstep.Implementation
	complete *completestep.Implementation
}

func (i *Implementation) PipelineFan(ctx context.Context, orders <-chan model.OrderId, workerId int, result chan model.PipelineOrder) {
	createCh := i.create.Pipeline(ctx, model.WorkerId(workerId), orders)

	const limit = 5
	fanOutProcess := make([]<-chan model.PipelineOrder, limit)
	for iter := 0; iter < limit; iter++ {
		fanOutProcess[iter] = i.process.Pipeline(ctx, createCh)
	}

	pipeline := i.complete.Pipeline(ctx, fanIn(ctx, fanOutProcess))
	for order := range pipeline {
		result <- order
		//fmt.Printf("Order: %d, storage: %d, pickup point: %d, worker: %v, tracking: %v \n", order.Order.Id, order.Order.Storage, order.Order.PickupPoint, order.Order.Worker, order.Order.Tracking)
		if order.Err != nil {
			log.Printf("error while processing order: [%v], err: [%v]", order.Order, order.Err)
		}
	}
}

func fanIn(ctx context.Context, chans []<-chan model.PipelineOrder) <-chan model.PipelineOrder {
	multiplexed := make(chan model.PipelineOrder)

	var wg sync.WaitGroup
	for _, ch := range chans {
		wg.Add(1)

		go func(ch <-chan model.PipelineOrder) {
			defer wg.Done()
			for v := range ch {
				select {
				case <-ctx.Done():
					return
				case multiplexed <- v:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(multiplexed)
	}()

	return multiplexed
}
