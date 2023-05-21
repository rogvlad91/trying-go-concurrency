package generator

import (
	"context"
	"trying-concurrency-go/internal/model"
)

func GenerateOrderIds(ctx context.Context, limit uint64) <-chan model.OrderId {
	result := make(chan model.OrderId)

	go func() {
		defer close(result)

		var counter uint64

		for {
			counter++

			if counter == limit {
				return
			}

			select {
			case <-ctx.Done():
				return
			case result <- model.OrderId(counter):
			}
		}
	}()

	return result
}
