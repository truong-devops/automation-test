Feature: Get orders
    As an API consumer
    I want to retrieve orders
    So that I can see existing orders

    Background:
        Given the following orders sxist:
            | alias   | customer_name | item     | quanlity | total_amount | status    |
            | order_a | An Nguyen     | Keyboard | 1        | 12000        | pending   |
            | order_b | Binh Tran     | Mouse    | 2        | 36000        | confirmed |
    Scenario: Get all existing orders
        When I send a Get request to "/api/v1/orders"
        Then ther response status should be 200
        And thew response field "count" should equal "2"