# Automation Test Microservices

Hai microservice Go nhỏ để thực hành API automation test, Godog, BDD, Gherkin và Cucumber.

## Kiến trúc

Mỗi service có cấu trúc độc lập:

```text
<service>/
├── cmd/server/main.go
├── internal/
│   ├── delivery/http/
│   ├── delivery/rpc/
│   ├── usecase/
│   ├── repository/
│   └── domain/
├── pkg/httputil/
├── configs/
├── migrations/
├── api/
├── Dockerfile
└── go.mod
```

- REST: CRUD để test qua HTTP.
- gRPC: `orders.GetOrder` và `payments.GetPaymentByOrderID`; gRPC reflection cũng được bật.
- MongoDB 7: hai database độc lập `orders_db` và `payments_db` trên cùng Mongo instance.
- Liên kết giữa service là tùy chọn. CRUD không yêu cầu service còn lại phải online.

## Chạy hệ thống

```bash
docker compose up --build -d
docker compose ps
```

Endpoints:

| Service | REST | gRPC |
|---|---|---|
| Orders | `localhost:8081` | `localhost:9091` |
| Payments | `localhost:8082` | `localhost:9092` |
| MongoDB | `localhost:27018` | - |

Tắt hệ thống (giữ dữ liệu):

```bash
docker compose down
```

Xóa luôn dữ liệu test:

```bash
docker compose down -v
```

## REST API

Orders:

```text
GET    /health
GET    /api/v1/orders
POST   /api/v1/orders
GET    /api/v1/orders/{id}
PUT    /api/v1/orders/{id}
DELETE /api/v1/orders/{id}
```

Payments:

```text
GET    /health
GET    /api/v1/payments
POST   /api/v1/payments
GET    /api/v1/payments/{id}
PUT    /api/v1/payments/{id}
DELETE /api/v1/payments/{id}
```

Tạo order:

```bash
curl -sS -X POST http://localhost:8081/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{
    "customer_name": "An Nguyen",
    "item": "Mechanical Keyboard",
    "quantity": 1,
    "total_amount": 1200000,
    "status": "pending"
  }'
```

Lấy `id` trong response làm `ORDER_ID`, sau đó tạo payment:

```bash
curl -sS -X POST http://localhost:8082/api/v1/payments \
  -H 'Content-Type: application/json' \
  -d "{
    \"order_id\": \"$ORDER_ID\",
    \"amount\": 1200000,
    \"method\": \"bank_transfer\",
    \"status\": \"completed\"
  }"
```

Hai request sau minh họa service gọi nhau qua gRPC:

```bash
curl -sS "http://localhost:8081/api/v1/orders/$ORDER_ID?include_payment=true"
curl -sS "http://localhost:8082/api/v1/payments/$PAYMENT_ID?include_order=true"
```

Giá trị hợp lệ:

- Order status: `pending`, `confirmed`, `cancelled`.
- Payment method: `card`, `bank_transfer`, `cash`.
- Payment status: `pending`, `completed`, `failed`, `refunded`.

## Kiểm tra code Go

```bash
go test ./orders/... ./payments/...
```

Các file `.proto` trong `api/` là hợp đồng gRPC. Implementation dùng trực tiếp Protobuf well-known types, nên build không phụ thuộc vào việc cài `protoc` trên máy.
