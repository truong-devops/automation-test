package automationtests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cucumber/godog"
)

func tableToMaps(table *godog.Table) ([]map[string]string, error) {

	headers := table.Rows[0].Cells

	var result []map[string]string

	for _, row := range table.Rows[1:] {
		data := make(map[string]string)

		for index, cell := range row.Cells {
			data[headers[index].Value] = cell.Value
		}

		result = append(result, data)
	}

	return result, nil
}

func buildOrderPayload(data map[string]string) ([]byte, error) {

	quantity, err := strconv.Atoi(data["quantity"])
	if err != nil {
		return nil, err
	}

	totalAmount, err := strconv.Atoi(data["total_amount"])

	if err != nil {
		return nil, err
	}

	payload := map[string]any{
		"customer_name": data["customer_name"],
		"item":          data["item"],
		"quantity":      quantity,
		"total_amount":  totalAmount,
		"status":        data["status"],
	}

	return json.Marshal(payload)
}

func (w *World) theFollowingOrdersExist(ctx context.Context, table *godog.Table) error {

	rows, err := tableToMaps(table)
	if err != nil {
		return err
	}

	for _, data := range rows {

		body, err := buildOrderPayload(data)
		if err != nil {
			return err
		}

		err = w.sendRequestWithBody( ctx, http.MethodPost, "/api/v1/orders", body, )

		if err != nil {
			return err
		}

		if w.lastStatus != http.StatusCreated {
			return fmt.Errorf( "failed to create fixture %s", data["alias"], )
		}

		var response struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		if err := json.Unmarshal( w.lastBody, &response, ); err != nil {
			return err
		}

		w.aliases[data["alias"]] = response.Data.ID
	}

	return nil
}

func (w *World) theFollowingOrderRequestFixtures(table *godog.Table) error {

	rows, err := tableToMaps(table)
	if err != nil {
		return err
	}

	for _, data := range rows {

		body, err := buildOrderPayload(data)
		if err != nil {
			return err
		}

		w.requestFixtures[data["alias"]] = body
	}

	return nil
}
