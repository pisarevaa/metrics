package server_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/pisarevaa/metrics/internal/server"
)

type ServerTestSuite struct {
	suite.Suite
	config server.Config
	logger *zap.SugaredLogger
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.config = server.GetConfig()
	suite.logger = server.GetLogger()
}

func TestAgentSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func testRequest(suite *ServerTestSuite, ts *httptest.Server, method, url, body string) (*http.Response, string) {
	bytesString := []byte(body)
	req, err := http.NewRequest(method, ts.URL+url, bytes.NewBuffer(bytesString))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := ts.Client().Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)
	return resp, string(respBody)
}

func (suite *ServerTestSuite) TestServerSaveLogs() {
	ts := httptest.NewServer(server.MetricsRouter(suite.config, suite.logger))
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
				json:       false,
				response:   "",
			},
			url:    "/update/gauge/HeapAlloc/1.25",
			body:   "",
			method: "POST",
		},
		{
			name: "add counter metric success",
			want: want{
				statusCode: 200,
				json:       false,
				response:   "",
			},
			url:    "/update/counter/PollCount/4",
			body:   "",
			method: "POST",
		},
		{
			name: "add wrong metric type",
			want: want{
				statusCode: 400,
				json:       false,
				response:   "",
			},
			url:    "/update/test/HeapAlloc/1.25",
			body:   "",
			method: "POST",
		},
		{
			name: "add empty metric value",
			want: want{
				statusCode: 404,
				json:       false,
				response:   "",
			},
			url:    "/update/counter/HeapAlloc/",
			body:   "",
			method: "POST",
		},
		{
			name: "add wrong metric value",
			want: want{
				statusCode: 400,
				json:       false,
				response:   "",
			},
			url:    "/update/counter/HeapAlloc/test",
			body:   "",
			method: "POST",
		},
		{
			name: "get gauge metric success",
			want: want{
				statusCode: 200,
				json:       false,
				response:   "1.25",
			},
			url:    "/value/gauge/HeapAlloc",
			body:   "",
			method: "GET",
		},
		{
			name: "get not found metric",
			want: want{
				statusCode: 404,
				json:       false,
				response:   "",
			},
			url:    "/value/gauge/NotFound",
			body:   "",
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
			body:   "",
			method: "GET",
		},
	}
	for _, tt := range tests {
		resp, body := testRequest(suite, ts, tt.method, tt.url, tt.body)
		defer resp.Body.Close()
		suite.Require().Equal(tt.want.statusCode, resp.StatusCode)
		if tt.want.response != "" {
			suite.Require().Equal(tt.want.response, body)
		}
	}
}

func (suite *ServerTestSuite) TestServerSaveLogsJSON() {
	ts := httptest.NewServer(server.MetricsRouter(suite.config, suite.logger))
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
				response:   `{"id":"HeapAlloc","type":"gauge","delta":0,"value":1.25}`,
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
				response:   `{"id":"PollCount","type":"counter","delta":4,"value":0}`,
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
			body:   `{"id":"HeapAlloc","type":"test","value":1.25,"delta":0}`,
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
				response:   `{"id":"HeapAlloc","type":"gauge","delta":0,"value":1.25}`,
			},
			url:    "/value/",
			body:   `{"id": "HeapAlloc", "type": "gauge"}`,
			method: "POST",
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
			method: "POST",
		},
		{
			name: "get all metrics success",
			want: want{
				statusCode: 200,
				json:       false,
				response:   "HeapAlloc: 1.25\nPollCount: 4\n",
			},
			url:    "/",
			body:   "",
			method: "GET",
		},
	}
	for _, tt := range tests {
		resp, body := testRequest(suite, ts, tt.method, tt.url, tt.body)
		defer resp.Body.Close()
		suite.Require().Equal(tt.want.statusCode, resp.StatusCode)
		if tt.want.response != "" {
			if tt.want.json {
				suite.Require().JSONEq(tt.want.response, body)
				suite.Require().Equal(tt.want.response, body)
			}
		}
	}
}
