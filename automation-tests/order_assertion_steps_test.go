package automationtests

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func (w *World) responseStatusShouldBe(expected int) error {

	if w.lastStatus != expected {
		return fmt.Errorf("expected status %d, got %d, body=%s", expected, w.lastStatus, string(w.lastBody))
	}
	return nil
}

func (w *World) responseFieldShouldEqual(fieldPath string, expected string) error {
	//Sử dụng thư viện gjson
	result := gjson.GetBytes(w.lastBody, fieldPath)
	if !result.Exists() {
		return fmt.Errorf("field %q not found", fieldPath)
	}

	actual := result.String()
	if actual != expected {
		return fmt.Errorf("field %q: expected %q, got %q", fieldPath, expected, actual)
	}

	return nil
}
