// Тесты по скорости сохранения и получения метрик.
package speed_test

import (
	"bytes"

	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/server/storage"
)

type ServerTestSuite struct {
	suite.Suite
	config server.Config
	logger *zap.SugaredLogger
}

func testRequest(suite *ServerTestSuite, ts *httptest.Server, method, url, body string) *http.Response {
	bytesString := []byte(body)
	req, err := http.NewRequest(method, ts.URL+url, bytes.NewBuffer(bytesString))
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "")
	resp, err := ts.Client().Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	return resp
}

func BenchmarkUpdateMetric(b *testing.B) {
	s := new(ServerTestSuite)
	s.SetT(&testing.T{})
	s.config = server.Config{
		Host:          "localhost:8080",
		StoreInterval: 300,
		Restore:       false,
	}
	s.logger = server.GetLogger()
	repo := storage.NewMemStorage()
	ts := httptest.NewServer(server.MetricsRouter(s.config, s.logger, repo))
	defer ts.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		resp := testRequest(s, ts, "POST", "/updates/", `[{"id": "HeapAlloc", "type": "gauge", "value": 1.25}]`)
		defer resp.Body.Close()
		b.StopTimer()
		s.Require().Equal(200, resp.StatusCode)
	}
}

func BenchmarkGetMetric(b *testing.B) {
	s := new(ServerTestSuite)
	s.SetT(&testing.T{})
	s.config = server.Config{
		Host:          "localhost:8080",
		StoreInterval: 300,
		Restore:       false,
	}
	s.logger = server.GetLogger()
	repo := storage.NewMemStorage()
	ts := httptest.NewServer(server.MetricsRouter(s.config, s.logger, repo))
	defer ts.Close()

	resp := testRequest(s, ts, "POST", "/updates/", `[{"id": "HeapAlloc", "type": "counter", "value": 2}]`)
	defer resp.Body.Close()
	s.Require().Equal(200, resp.StatusCode)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		resp = testRequest(s, ts, "POST", "/value/", `{"id": "HeapAlloc", "type": "counter"}`)
		defer resp.Body.Close()
		b.StopTimer()
		s.Require().Equal(200, resp.StatusCode)
	}
}
