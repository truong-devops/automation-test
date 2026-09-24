Bạn chạy:
go test -v
   │
   ▼
Go tìm các file *_test.go
   │
   ▼
suite_test.go
   │
   └── TestFeatures(t *testing.T)
          │
          ▼
      godog.TestSuite
          │
          ├── Paths: features/
          ├── Tags: GODOG_TAGS
          ├── ScenarioInitializer: InitializeScenario
          └── TestingT: t
          │
          ▼
      suite.Run()
          │
          ▼
Godog đọc tất cả file .feature
          │
          ├── create_order.feature
          ├── get_orders.feature
          ├── update_order.feature
          └── delete_order.feature
          │
          ▼
Godog lấy từng Scenario
          │
          ▼
order_steps_test.go
InitializeScenario()
          │
          ├── tạo World
          ├── đăng ký Before / After
          ├── map Given
          ├── map When
          └── map Then
          │
          ▼
      chạy Scenario