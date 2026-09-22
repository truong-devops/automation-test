@orders @get
Feature: Get orders
    As an API consumer
    I want to retrieve orders
    So that I can view order data

    Scenario: GET an order by ID
        Given the following orders exist:
            | alias   | customer_name | item  | quantity | total_amount | status    |
            | order_c | Ci Ci         | Nokia | 3        | 33000        | confirmed |

        When I send a GET request to "/api/v1/orders/{order_c}"

        Then the response status should be 200
        And the response field "data.customer_name" should equal "Ci Ci"
        And the response field "data.item" should equal "Nokia"
        And the response field "data.total_amount" should equal "33000"
        And the response field "data.status" should equal "confirmed"


    Scenario: GET all orders
        Given the following orders exist:
            | alias   | customer_name | item    | quantity | total_amount | status    |
            | order_a | An An         | Iphone  | 1        | 11000        | confirmed |
            | order_b | Bin Bin       | Samsung | 2        | 22000        | pending   |

        When I send a GET request to "/api/v1/orders"

        Then the response status should be 200
        And the response field "count" should equal "2"