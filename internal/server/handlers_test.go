package server_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pisarevaa/metrics/internal/server"
)

func testRequest(t *testing.T, ts *httptest.Server, method, url string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+url, nil)
	require.NoError(t, err)
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, string(respBody)
}

func TestServerSaveLogs(t *testing.T) {
	config := server.GetConfig()
	ts := httptest.NewServer(server.MetricsRouter(config))
	defer ts.Close()

	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name   string
		url    string
		method string
		want   want
	}{
		{
			name: "add gauge metric success",
			want: want{
				statusCode: 200,
				response:   "",
			},
			url:    "/update/gauge/HeapAlloc/1.25",
			method: "POST",
		},
		{
			name: "add counter metric success",
			want: want{
				statusCode: 200,
				response:   "",
			},
			url:    "/update/counter/PollCount/4",
			method: "POST",
		},
		{
			name: "add wrong metric type",
			want: want{
				statusCode: 400,
				response:   "",
			},
			url:    "/update/test/HeapAlloc/1.25",
			method: "POST",
		},
		{
			name: "add empty metric value",
			want: want{
				statusCode: 404,
				response:   "",
			},
			url:    "/update/counter/HeapAlloc/",
			method: "POST",
		},
		{
			name: "add wrong metric value",
			want: want{
				statusCode: 400,
				response:   "",
			},
			url:    "/update/counter/HeapAlloc/test",
			method: "POST",
		},
		{
			name: "get gauge metric success",
			want: want{
				statusCode: 200,
				response:   "1.25",
			},
			url:    "/value/gauge/HeapAlloc",
			method: "GET",
		},
		{
			name: "get not found metric",
			want: want{
				statusCode: 404,
				response:   "",
			},
			url:    "/value/gauge/NotFound",
			method: "GET",
		},
		{
			name: "get all metrics success",
			want: want{
				statusCode: 200,
				response:   "HeapAlloc: 1.25\nPollCount: 4\n",
			},
			url:    "/",
			method: "GET",
		},
	}
	for _, tt := range tests {
		resp, body := testRequest(t, ts, tt.method, tt.url)
		defer resp.Body.Close()
		assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		if tt.want.response != "" {
			assert.Equal(t, tt.want.response, body)
		}
	}
}
