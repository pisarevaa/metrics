package server_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pisarevaa/metrics/internal/server"
)

func testRequest(t *testing.T, ts *httptest.Server, method, url, body string) (*http.Response, string) {
	bytesString := []byte(body)
	req, err := http.NewRequest(method, ts.URL+url, bytes.NewBuffer(bytesString))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, string(respBody)
}

func TestServerSaveLogs(t *testing.T) {
	config := server.GetConfig()
	logger := server.GetLogger()
	ts := httptest.NewServer(server.MetricsRouter(config, logger))
	defer ts.Close()

	type want struct {
		statusCode int
		json       bool
		response   string
	}
	tests := []struct {
		name   string
		url    string
		method string
		body   string
		want   want
	}{
		{
			name: "add gauge metric success",
			want: want{
				statusCode: 200,
				json:       true,
				response:   `{"id": "HeapAlloc", "type": "gauge", "value": 1.25, "delta": 0}`,
			},
			url:    "/update/",
			body:   `{"id": "HeapAlloc", "type": "gauge", "value": 1.25}`,
			method: "POST",
		},
		{
			name: "add counter metric success",
			want: want{
				statusCode: 200,
				json:       true,
				response:   `{"id": "PollCount", "type": "counter", "value": 0, "delta": 4}`,
			},
			url:    "/update/",
			body:   `{"id": "PollCount", "type": "counter", "delta": 4}`,
			method: "POST",
		},
		{
			name: "add wrong metric type",
			want: want{
				statusCode: 400,
				json:       true,
				response:   "",
			},
			url:    "/update/",
			body:   `{"id": "HeapAlloc", "type": "test", "value": 1.25}`,
			method: "POST",
		},
		{
			name: "add empty metric value",
			want: want{
				statusCode: 400,
				json:       true,
				response:   "",
			},
			url:    "/update/",
			body:   `{"id": "HeapAlloc", "type": "counter"}`,
			method: "POST",
		},
		{
			name: "add wrong metric value",
			want: want{
				statusCode: 400,
				json:       true,
				response:   "",
			},
			url:    "/update/",
			body:   `{"id": "HeapAlloc", "type": "counter", "delta": "test"}`,
			method: "POST",
		},
		{
			name: "get gauge metric success",
			want: want{
				statusCode: 200,
				json:       true,
				response:   `{"id": "HeapAlloc", "type": "gauge", "value": 1.25, "delta": 0}`,
			},
			url:    "/value/",
			body:   `{"id": "HeapAlloc", "type": "gauge"}`,
			method: "GET",
		},
		{
			name: "get not found metric",
			want: want{
				statusCode: 404,
				json:       true,
				response:   "",
			},
			url:    "/value/",
			body:   `{"id": "NotFound", "type": "gauge"}`,
			method: "GET",
		},
		{
			name: "get all metrics success",
			want: want{
				statusCode: 200,
				json:       false,
				response:   "HeapAlloc: 1.25\nPollCount: 4\n",
			},
			url:    "/",
			body:   ``,
			method: "GET",
		},
	}
	for _, tt := range tests {
		resp, body := testRequest(t, ts, tt.method, tt.url, tt.body)
		defer resp.Body.Close()
		assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		if tt.want.response != "" {
			if tt.want.json {
				assert.JSONEq(t, tt.want.response, body)
			} else {
				assert.Equal(t, tt.want.response, body)
			}
		}
	}
}
