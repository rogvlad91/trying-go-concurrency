package main

import (
	"context"
	"encoding/json"
	"fmt"
	"homework/internal/model"
	"homework/internal/pkg/generator"
	"homework/internal/pkg/service/gateway"
	completestep "homework/internal/pkg/service/gateway/steps/complete"
	createstep "homework/internal/pkg/service/gateway/steps/create"
	processstep "homework/internal/pkg/service/gateway/steps/process"
	"log"
	"sync"
	"time"
)

const (
	workerTotal = 2
	orderTotal  = 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ids := generator.GenerateOrderIds(ctx, orderTotal)

	create := createstep.New()
	process := processstep.New()
	complete := completestep.New()

	server := gateway.New(create, process, complete)

	start := time.Now().UTC()

	result := make(chan model.PipelineOrder)

	var wg sync.WaitGroup
	for i := 0; i < workerTotal; i++ {
		wg.Add(1)
		go worker(ctx, &wg, server, ids, i, result)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for order := range result {
		data, err := json.Marshal(order)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(data))

		if order.Err != nil {
			log.Printf("error while processing order: [%v], err: [%v]", order.Order, order.Err)
		}
	}

	wg.Wait()

	fmt.Printf("Total duration: %f", time.Since(start).Seconds())
}

func worker(ctx context.Context, wg *sync.WaitGroup, server *gateway.Implementation, orders <-chan model.OrderId, workerId int, result chan model.PipelineOrder) {
	defer wg.Done()
	server.PipelineFan(ctx, orders, workerId, result)
}
