package internal

import (
	"context"
	"fmt"
	"github.com/hwebz/simple-go-grpc-gateway/protogen/golang/orders"
	"log"
)

// OrderService should implement the OrdersServer interface generated from grpc.
//
// UnimplementedOrderServer must be embedded to have forwarded compatible implementations.
type OrderService struct {
	db *DB
	orders.UnimplementedOrdersServer
}

// NewOrderService creates a new OrderService
func NewOrderService(db *DB) OrderService {
	return OrderService{db: db}
}

// AddOrder implements the AddOrder method of the grpc OrdersServer interface to add a new order
func (o *OrderService) AddOrder(_ context.Context, req *orders.PayloadWithSingleOrder) (*orders.Empty, error) {
	log.Printf("Received an add-order request")

	err := o.db.AddOrder(req.GetOrder())

	return &orders.Empty{}, err
}

// UpdateOrder implements the UpdateOrder method of the grpc OrdersServer interface to update an order
func (o *OrderService) UpdateOrder(_ context.Context, req *orders.PayloadWithSingleOrder) (*orders.Empty, error) {
	log.Printf("Received an update order request")

	o.db.UpdateOrder(req.GetOrder())

	return &orders.Empty{}, nil
}

// GetOrder implements the GetOrder method of the grpc OrdersServer interface to fetch an order for a given orderID
func (o *OrderService) GetOrder(_ context.Context, req *orders.PayloadWithOrderID) (*orders.PayloadWithSingleOrder, error) {
	log.Printf("Received get order request with orderID: %d", req.GetOrderId())

	order := o.db.GetOrderByID(req.GetOrderId())
	if order == nil {
		return nil, fmt.Errorf("Order not found for the orderID: %d", req.GetOrderId())
	}

	return &orders.PayloadWithSingleOrder{Order: order}, nil
}

// RemoveOrder implements the RemoveOrder method of the grpc OrdersServer interface to remove an order
func (o *OrderService) RemoveOrder(_ context.Context, req *orders.PayloadWithOrderID) (*orders.Empty, error) {
	log.Printf("Received remove order request")

	o.db.RemoveOrder(req.GetOrderId())

	return &orders.Empty{}, nil
}
