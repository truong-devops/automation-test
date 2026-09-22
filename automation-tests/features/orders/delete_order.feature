@orders @delete
Feature: Delete order
    As an API consumer
    I want to delete orders
    So that unwanted order data can be removed

    Scenario: Delete an existing order successfully
        Given the following orders exist:
            | alias   | customer_name | item     | quantity | total_amount | status    |
            | order_e | En En         | Eyeclass | 1        | 3000         | confirmed |

        When I send a DELETE request to "/api/v1/orders/{order_e}"
        Then the response status should be 204

        When I send a GET request to "/api/v1/orders/{order_e}"
        Then the response status should be 404