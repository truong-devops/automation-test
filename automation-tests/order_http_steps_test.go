package automationtests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
//	"strings"
	"regexp"
)

// func (w *World) resolveAliases(path string) (string, error) {

// 	result := path

// 	for alias, id := range w.aliases {
// 		result = strings.ReplaceAll(result, "{"+alias+"}", id)
// 	}

// 	if strings.Contains(result, "{") {
// 		return "", fmt.Errorf("unknown alias in %q", result)
// 	}

// 	return result, nil
// }

var aliasRegex = regexp.MustCompile(`\{([^}]+)\}`)

func (w *World) resolveAliases(path string) (string, error) {
	var err error

	result := aliasRegex.ReplaceAllStringFunc(path, func(match string) string {
		aliasName := match[1 : len(match)-1]

		if id, exists := w.aliases[aliasName]; exists {
			return id
		}

		err = fmt.Errorf("unknown alias %q", aliasName)
		return match
	})

	if err != nil {
		return "", fmt.Errorf("unknown alias in %q", result)
	}

	return result, nil
}

func (w *World) executeRequest(request *http.Request) error {

	response, err := w.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	w.lastStatus = response.StatusCode
	w.lastBody = body

	fmt.Println("Response status:", w.lastStatus)

	fmt.Println("Response body:", string(w.lastBody))

	return nil
}

func (w *World) sendRequestUsingFixture(ctx context.Context, method string, path string, fixtureAlias string) error {

	body, exists := w.requestFixtures[fixtureAlias]

	if !exists {
		return fmt.Errorf("fixture %q not found", fixtureAlias)
	}

	return w.sendRequestWithBody(ctx, method, path, body)
}

func (w *World) sendRequestWithBody(ctx context.Context, method string, path string, body []byte) error {

	resolvedPath, err := w.resolveAliases(path)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, method, w.baseURL+resolvedPath, bytes.NewReader(body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	return w.executeRequest(request)
}

func (w *World) sendRequestWithoutBody(ctx context.Context, method string, path string) error {

	resolvedPath, err := w.resolveAliases(path)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, method, w.baseURL+resolvedPath, nil)
	if err != nil {
		return err
	}

	return w.executeRequest(request)
}
