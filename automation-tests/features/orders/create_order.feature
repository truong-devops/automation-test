@orders @post
Feature: Create order
    As an API consumer
    I want to create orders
    So that new order data can be stored

    Scenario: Create a new order successfully
        Given the following order request fixtures:
            | alias   | customer_name | item     | quantity | total_amount | status  |
            | order_a | An Nguyen     | Keyboard | 1        | 12000        | pending |

        When I send a POST request to "/api/v1/orders" using fixture "order_a"

        Then the response status should be 201
        And the response field "data.customer_name" should equal "An Nguyen"
        And the response field "data.item" should equal "Keyboard"
        And the response field "data.quantity" should equal "1"
        And the response field "data.total_amount" should equal "12000"
        And the response field "data.status" should equal "pending"