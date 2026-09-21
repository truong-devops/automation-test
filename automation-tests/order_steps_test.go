package automationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func InitializeScenario(scenario *godog.ScenarioContext) {
	world := newWorld()

	scenario.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		err := world.connectMongo(ctx)
		if err != nil {
			return ctx, err
		}

		err = world.resetMongo(ctx)
		if err != nil {
			return ctx, err
		}

		return ctx, nil
	})

	scenario.After(func(ctx context.Context, sc *godog.Scenario, scenarioErr error) (context.Context, error) {
		err := world.closeMongo(ctx)

		return ctx, err
	})

	scenario.When(
		`^I send a (POST|PUT) request to "([^"]*)" with body:$`, world.sendRequestWithBody,
	)

	scenario.Given(
		`^the following orders exist:$`, world.theFollowingOrdersExist,
	)

	scenario.When(
		`^I send a (GET|DELETE) request to "([^"]*)"$`, world.sendRequestWithoutBody,
	)

	scenario.Then(
		`^the response status should be (\d+)$`, world.responseStatusShouldBe,
	)

	scenario.Then(
		`^the response field "([^"]*)" should equal "([^"]*)"$`, world.responseFieldShouldEqual,
	)
}

func (w *World) theFollowingOrdersExist(ctx context.Context, table *godog.Table) error {
	headers := table.Rows[0].Cells

	for _, row := range table.Rows[1:] {
		data := make(map[string]string)
		for index, cell := range row.Cells {
			columnName := headers[index].Value
			data[columnName] = cell.Value
		}

		quantity, err := strconv.Atoi(data["quantity"])
		if err != nil {
			return err
		}

		totalAmount, err := strconv.Atoi(data["total_amount"])
		if err != nil {
			return err
		}

		payload := map[string]any{
			"customer_name": data["customer_name"],
			"item":          data["item"],
			"quantity":      quantity,
			"total_amount":  totalAmount,
			"status":        data["status"],
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		docString := &godog.DocString{
			Content: string(body),
		}

		err = w.sendRequestWithBody(ctx, http.MethodPost, "/api/v1/orders", docString)
		if err != nil {
			return err
		}

		if w.lastStatus != http.StatusCreated {
			return fmt.Errorf(
				"failed to create fixture %s: status %d, body: %s",
				data["alias"],
				w.lastStatus,
				string(w.lastBody),
			)
		}

		fmt.Printf(
			"Fixture %s created successfully!\n", data["alias"],
		)
	}
	return nil
}

func (w *World) sendRequestWithBody(
	ctx context.Context,
	method string,
	path string,
	body *godog.DocString,
) error {
	fmt.Println("Method:", method)
	fmt.Println("Path:", path)
	fmt.Println("Body:", body.Content)

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		w.baseURL+path,
		strings.NewReader(body.Content),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := w.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	w.lastStatus = response.StatusCode
	w.lastBody = responseBody

	fmt.Println("Response status:", w.lastStatus)
	fmt.Println("Respone body", string(w.lastBody))

	return nil
}

func (w *World) sendRequestWithoutBody(
	ctx context.Context,
	method string,
	path string,
) error {
	request, err := http.NewRequestWithContext(ctx, method, w.baseURL+path, nil)
	if err != nil {
		return err
	}

	response, err := w.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	w.lastStatus = response.StatusCode
	w.lastBody = responseBody

	fmt.Println("Response status:", w.lastStatus)
	fmt.Println("Respone body:", string(w.lastBody))

	return nil
}

func (w *World) responseStatusShouldBe(expected int) error {
	if w.lastStatus != expected {
		return fmt.Errorf(
			"expected status %d, got %d",
			expected,
			w.lastStatus,
		)
	}

	return nil
}

func (w *World) responseFieldShouldEqual(
	fieldPath string,
	expected string,
) error {
	var response map[string]any

	err := json.Unmarshal(w.lastBody, &response)
	if err != nil {
		return err
	}

	parts := strings.Split(fieldPath, ".")

	var current any = response

	for _, part := range parts {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return fmt.Errorf(
				"field %q not an object",
				part,
			)
		}

		value, exists := currentMap[part]
		if !exists {
			return fmt.Errorf(
				"field %q not found",
				part,
			)
		}

		current = value
	}

	actual := fmt.Sprint(current)

	if actual != expected {
		return fmt.Errorf(
			"field %q: expected %q, got %q", fieldPath, expected, actual,
		)
	}

	return nil
}
