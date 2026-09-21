Feature: Orders API
    As an API consumer
    I want to manage orders
    So that I can work with order data
    
    Scenario: Create a new order
        When I send a POST request to "/api/v1/orders" with body:
            """
            {
                "customer_name": "An Nguyen",
                "item": "Keyboard",
                "quantity": 1,
                "total_amount": 12000,
                "status": "pending"
            }
            """
        Then the response status should be 201
        And the response field "data.customer_name" should equal "An Nguyen"
        And the response field "data.status" should equal "pending"

    Scenario: GET all orders
        Given the following orders exist:
            | alias   | customer_name | item   | quantity | total_amount | status    |
            | order_a | An An         | Iphone | 1        | 11000        | confirmed |
            | order_b | Bin Bin       | Samsung| 2        | 22000        | pending   |
        When I send a GET request to "/api/v1/orders"
        Then the response status should be 200
        And the response field "count" should equal "2"