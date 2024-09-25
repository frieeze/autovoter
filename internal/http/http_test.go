package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	// Create a test server that returns a 200 OK response with a JSON body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"key": "value"}`)
	}))
	defer ts.Close()

	var (
		body = struct {
			Value string `json:"value"`
		}{
			Value: "test",
		}

		recipient struct {
			Key string `json:"key"`
		}
	)

	err := Post(context.Background(), ts.URL, body, &recipient)
	assert.NoError(t, err)
	assert.Equal(t, "value", recipient.Key)
}

func TestPost_error(t *testing.T) {
	// Create a test server that returns a 500 Internal Server Error response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	var (
		body = struct {
			Value string `json:"value"`
		}{
			Value: "test",
		}

		recipient struct {
			Key string `json:"key"`
		}
	)
	err := Post(context.Background(), ts.URL, body, recipient)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrRequestFailed)
}
