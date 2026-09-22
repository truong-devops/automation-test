@orders @put
Feature: Update order
    As an API consumer
    I want to update existing orders
    So that order information can be changed

    Scenario: Update an existing order successfully
        Given the following orders exist:
            | alias   | customer_name | item   | quantity | total_amount | status    |
            | order_d | Dan Dan       | Laptop | 1        | 15000        | confirmed |

        And the following order request fixtures:
            | alias          | customer_name   | item   | quantity | total_amount | status  |
            | update_order_d | Dan Dan updated | new PC | 1        | 36000        | pending |

        When I send a PUT request to "/api/v1/orders/{order_d}" using fixture "update_order_d"

        Then the response status should be 200
        And the response field "data.customer_name" should equal "Dan Dan updated"
        And the response field "data.item" should equal "new PC"
        And the response field "data.quantity" should equal "1"
        And the response field "data.total_amount" should equal "36000"
        And the response field "data.status" should equal "pending"