package agent

import (
	"github.com/go-resty/resty/v2"
)

func NewClient() *resty.Client {
	client := resty.New()
	return client
}
