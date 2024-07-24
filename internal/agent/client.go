package agent

import (
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const retries = 3
const retryWaitTime = 1
const retryMaxWaitTime = 20

// Создание клиента для оптравки метрик.
// Ключевой момент что сервер может быть недостпуен и ниже реализована логика таймаута.
func NewClient() *resty.Client {
	client := resty.New()
	client.
		SetRetryCount(retries).
		SetRetryWaitTime(retryWaitTime * time.Second).
		SetRetryMaxWaitTime(retryMaxWaitTime * time.Second).
		AddRetryCondition(
			func(_ *resty.Response, err error) bool {
				if err == nil {
					return false
				}
				return strings.Contains(err.Error(), "connect: connection refused")
			},
		)
	return client
}
