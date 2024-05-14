package server_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/golang/mock/gomock"

	"github.com/pisarevaa/metrics/internal/server"
	mock "github.com/pisarevaa/metrics/internal/server/mocks"
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
	req.Header.Set("Accept-Encoding", "")
	resp, err := ts.Client().Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)
	return resp, string(respBody)
}

func testRequestWithGZIP(
	suite *ServerTestSuite, ts *httptest.Server, method, url, body, contentEncoding, acceptEncoding string,
) (*http.Response, string) {
	bodyBytes := []byte(body)
	var reqBody io.Reader
	if contentEncoding == "gzip" {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write(bodyBytes)
		suite.Require().NoError(err)
		err = zb.Close()
		suite.Require().NoError(err)
		reqBody = buf
	} else {
		reqBody = bytes.NewBuffer(bodyBytes)
	}
	req, err := http.NewRequest(method, ts.URL+url, reqBody)
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", contentEncoding)
	req.Header.Set("Accept-Encoding", acceptEncoding)
	resp, err := ts.Client().Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	if acceptEncoding == "gzip" {
		zr, errReader := gzip.NewReader(resp.Body)
		suite.Require().NoError(errReader)
		b, errIo := io.ReadAll(zr)
		suite.Require().NoError(errIo)
		return resp, string(b)
	}
	respBody, err := io.ReadAll(resp.Body)
	suite.Require().NoError(err)
	return resp, string(respBody)
}

func (suite *ServerTestSuite) TestServerUpdateAndGetMetrics() {
	storage := server.NewMemStorageRepo()
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockMetricsModel(ctrl)

	m.EXPECT().
		IsExist().
		Return(false).
		MaxTimes(10)
	ts := httptest.NewServer(server.MetricsRouter(suite.config, suite.logger, storage, m))
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

func (suite *ServerTestSuite) TestServerUpdateAndGetMetricsJSON() {
	storage := server.NewMemStorageRepo()
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockMetricsModel(ctrl)

	m.EXPECT().
		IsExist().
		Return(false).
		MaxTimes(10)
	ts := httptest.NewServer(server.MetricsRouter(suite.config, suite.logger, storage, m))
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
				response:   `{"id":"HeapAlloc","type":"gauge","value":1.25}`,
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

func (suite *ServerTestSuite) TestServerUpdateAndGetMetricsWithGZIP() {
	storage := server.NewMemStorageRepo()
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockMetricsModel(ctrl)

	m.EXPECT().
		IsExist().
		Return(false).
		MaxTimes(10)

	ts := httptest.NewServer(server.MetricsRouter(suite.config, suite.logger, storage, m))
	defer ts.Close()

	type want struct {
		statusCode int
		json       bool
		response   string
	}
	tests := []struct {
		name            string
		url             string
		method          string
		body            string
		contentEncoding string
		acceptEncoding  string
		want            want
	}{
		{
			name: "add gauge metric with gzip body",
			want: want{
				statusCode: 200,
				json:       true,
				response:   `{"id":"HeapAlloc","type":"gauge","delta":0,"value":1.25}`,
			},
			url:             "/update/",
			body:            `{"id": "HeapAlloc", "type": "gauge", "value": 1.25}`,
			method:          "POST",
			contentEncoding: "gzip",
			acceptEncoding:  "",
		},
		{
			name: "add gauge metric with gzip return",
			want: want{
				statusCode: 200,
				json:       true,
				response:   `{"id":"HeapAlloc","type":"gauge","delta":0,"value":1.25}`,
			},
			url:             "/update/",
			body:            `{"id": "HeapAlloc", "type": "gauge", "value": 1.25}`,
			method:          "POST",
			contentEncoding: "",
			acceptEncoding:  "gzip",
		},
		{
			name: "add gauge metric with gzip body and return",
			want: want{
				statusCode: 200,
				json:       true,
				response:   `{"id":"HeapAlloc","type":"gauge","delta":0,"value":1.25}`,
			},
			url:             "/update/",
			body:            `{"id": "HeapAlloc", "type": "gauge", "value": 1.25}`,
			method:          "POST",
			contentEncoding: "gzip",
			acceptEncoding:  "gzip",
		},
		{
			name: "get gauge metric with gzip success",
			want: want{
				statusCode: 200,
				json:       true,
				response:   `{"id":"HeapAlloc","type":"gauge","value":1.25}`,
			},
			url:             "/value/",
			body:            `{"id": "HeapAlloc", "type": "gauge"}`,
			method:          "POST",
			contentEncoding: "",
			acceptEncoding:  "gzip",
		},
	}
	for _, tt := range tests {
		resp, body := testRequestWithGZIP(suite, ts, tt.method, tt.url, tt.body, tt.contentEncoding, tt.acceptEncoding)
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

func (suite *ServerTestSuite) TestPing() {
	storage := server.NewMemStorageRepo()
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockMetricsModel(ctrl)

	m.EXPECT().
		Ping(gomock.Any()).
		Return(nil)
	m.EXPECT().
		IsExist().
		Return(true).
		MaxTimes(2)
	m.EXPECT().
		RestoreMetricsFromDB(gomock.Any()).
		Return(nil)

	ts := httptest.NewServer(server.MetricsRouter(suite.config, suite.logger, storage, m))
	defer ts.Close()

	resp, _ := testRequest(suite, ts, "GET", "/ping", "")
	defer resp.Body.Close()
	suite.Require().Equal(200, resp.StatusCode)
}

func (suite *ServerTestSuite) TestServerUpdateAndGetMetricsJSONBatch() {
	storage := server.NewMemStorageRepo()
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	m := mock.NewMockMetricsModel(ctrl)

	m.EXPECT().
		InsertRowsIntoDB(gomock.Any(), gomock.Any()).
		Return(nil)
	m.EXPECT().
		IsExist().
		Return(true).
		MaxTimes(2)
	m.EXPECT().
		RestoreMetricsFromDB(gomock.Any()).
		Return(nil)

	ts := httptest.NewServer(server.MetricsRouter(suite.config, suite.logger, storage, m))
	defer ts.Close()

	resp, _ := testRequest(suite, ts, "POST", "/updates/", `[{"id": "HeapAlloc", "type": "gauge", "value": 1.25}]`)
	defer resp.Body.Close()
	suite.Require().Equal(200, resp.StatusCode)
}
