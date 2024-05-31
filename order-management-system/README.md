# Order Management System - Microservices

## Tech Stack
```
1. Go 1.22+
2. Golang cosmtrek/air for hot-reloading
3. gRPC for communmication between services
4. RabbitMQ as message broker
5. Docker with Docker Compose
6. MongoDB as storage layer
7. Jaeger for service tracing
8. HashiCorp's Consul for service discovery
9. Stripe for payments
```

## Architectures
<img src="./screenshots/oms.jpg" alt="OMS" />

```
Orders Service
  -, Validate order details -> Talk with stock service
  -, CRUD of Orders
  -, Initiates the Payment Flow -> by sending an event

Stock Service
  -, Handles stock
  -, Validate order quantities
  -, Might return items as menu

Menu Service
  -, Stores items as menu

Payment Service
  -, Initiates a payment with a 3rd party provider (Stripe)
  -, Produces an order Paid/Cancelled event to orders/stock/kitchen

Kitchen Service
  -, Long running process of a "Simulated kitchen staff"
```

## Technical Guide 
> cd common && go mod init github.com/hwebz/oms-commons
> cd gateway && go mod init github.com/hwebz/oms-gateway
> cd kitchen && go mod init github.com/hwebz/oms-kitchen
> cd orders && go mod init github.com/hwebz/oms-orders
> cd payments && go mod init github.com/hwebz/oms-payments
> cd stock && go mod init github.com/hwebz/oms-stock
> go work init ./common ./gateway ./kitchen ./orders ./payments ./stock