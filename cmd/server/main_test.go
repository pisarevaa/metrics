package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/pisarevaa/metrics/internal/server"
)

func TestServerSaveLogs(t *testing.T) {
    type want struct {
        contentType string
        statusCode  int
    }
    tests := []struct {
        name    string
        request string
        want    want
    }{
        {
            name: "gauge success test",
            want: want{
                contentType: "text/plain",
                statusCode:  200,
            },
            request: "/update/gauge/HeapAlloc/1.25",
        },
        {
            name: "counter success test",
            want: want{
                contentType: "text/plain",
                statusCode:  200,
            },
            request: "/update/counter/PollCount/4",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            request := httptest.NewRequest(http.MethodPost, tt.request, nil)
            w := httptest.NewRecorder()
            storage := server.MemStorage{Metrics: server.MetricGroup{Gauge: make(map[string]float64), Counter: make(map[string]int64)}}
            h := http.HandlerFunc(storage.HandleMetrics)
            h(w, request)
            result := w.Result()
            assert.Equal(t, tt.want.statusCode, result.StatusCode)
            assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
            defer result.Body.Close()
        })
    }
}