package automationtests

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (w *World) responseStatusShouldBe(expected int) error {

	if w.lastStatus != expected {
		return fmt.Errorf("expected status %d, got %d, body=%s", expected, w.lastStatus, string(w.lastBody))
	}

	return nil
}

func (w *World) responseFieldShouldEqual(fieldPath string, expected string) error {

	var response map[string]any

	if err := json.Unmarshal(w.lastBody, &response); err != nil {
		return err
	}

	var current any = response

	for _, part := range strings.Split(fieldPath, ".") {

		currentMap, ok := current.(map[string]any)
		if !ok {
			return fmt.Errorf("field %q is not an object", part)
		}

		value, exists := currentMap[part]
		if !exists {
			return fmt.Errorf("field %q not found", part)
		}

		current = value
	}

	actual := fmt.Sprint(current)
	if actual != expected {
		return fmt.Errorf("field %q: expected %q, got %q", fieldPath, expected, actual)
	}

	return nil
}
