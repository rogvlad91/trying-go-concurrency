package model

import "time"

type OrderState string

const (
	OrderStateCreated    = OrderState("created")
	OrderStateProcessing = OrderState("processing")
	OrderStateCompleted  = OrderState("completed")
)

type OrderId uint64
type StorageId uint64
type PickupPointId uint64
type WorkerId uint64

type Order struct {
	Id          OrderId
	Storage     StorageId
	PickupPoint PickupPointId
	Worker      WorkerId
	Tracking    []OrderTracking
}

type OrderTracking struct {
	State OrderState
	Start time.Time
}

type PipelineOrder struct {
	Order Order
	Err   error
}
