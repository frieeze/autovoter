package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	// ErrRequestFailed is returned when a request status code is over 299
	ErrRequestFailed = fmt.Errorf("request failed")
)

// Post performs a POST request to the given route with data and unmarshal the response to recipient
func Post(ctx context.Context, route string, data interface{}, recipient interface{}) error {
	query, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed parsing request: %w", err)
	}

	fmt.Println("query", string(query))

	req, err := http.NewRequestWithContext(ctx, "POST", route, bytes.NewBuffer(query))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return fmt.Errorf("%w : %s", ErrRequestFailed, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("fail to decode response: %w", err)
	}

	err = json.Unmarshal(body, recipient)
	if err != nil {
		return fmt.Errorf("fail to unmarshal response: %w", err)
	}
	return nil
}
