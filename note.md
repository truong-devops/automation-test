Bạn chạy:
go test -v
    │
    ▼
suite_test.go
    │
    │ khởi động Godog
    ▼
features/orders.feature
    │
    │ đọc Given / When / Then
    ▼
InitializeScenario()
    │
    │ tìm step definition tương ứng
    ▼
order_steps_test.go
    │
    ├── Given → chuẩn bị dữ liệu test
    ├── When  → gọi API
    └── Then  → kiểm tra kết quả
            │
            ▼
          World
   lưu response/status/ID
            │
            ▼
     Orders Service :8081
            │
            ▼
          MongoDB