package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/pisarevaa/metrics/internal/server"
	"github.com/pisarevaa/metrics/internal/storage"
	"net/http"
)

func MetricsRouter() chi.Router {
	storage := storage.MemStorage{}
	storage.Init()
	server := server.Server{Storage: &storage}
    r := chi.NewRouter()
    r.Post("/update/{metricType}/{metricName}/{metricValue}", server.StoreMetrics)
	r.Get("/value/{metricType}/{metricName}", server.GetMetric)
	r.Get("/", server.GetAllMetrics)
    return r
}

// http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ> он возвращал текущее значение метрики в текстовом виде со статусом http.StatusOK.
// При попытке запроса неизвестной метрики сервер должен возвращать http.StatusNotFound.
// По запросу GET http://<АДРЕС_СЕРВЕРА>/ сервер должен отдавать HTML-страницу со списком имён и значений всех известных ему на текущий момент метрик.

func main() {
	err := http.ListenAndServe(":8080", MetricsRouter())
	if err != nil {
		panic(err)
	}
}
