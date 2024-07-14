.PHONY: help unit_test integration_test e2e_test test lint coverage_report

help:
	cat Makefile

run_agent:
	go run ./cmd/agent

run_server:
	go run ./cmd/server

test_agent:
	go clean -testcache && go test -v ./internal/agent

test_server:
	go clean -testcache && go test -v ./internal/server

pprof_heap_agent:
	go tool pprof -http=":9090" -seconds=9 http://localhost:8085/debug/pprof/profile

pprof_heap_agent:
	go tool pprof -http=":9090" -seconds=9 http://localhost:8085/debug/pprof/heap

pprof_heap_server:
	go tool pprof -http=":9090" -seconds=9 http://localhost:8080/debug/pprof/profile

pprof_heap_server:
	go tool pprof -http=":9090" -seconds=9 http://localhost:8080/debug/pprof/heap

test_server_speed:
	go clean -testcache && go test -bench=. -v ./internal/speed | grep "Benchmark"

mock_generate:
	mockgen -source=internal/server/storage/types.go -destination=internal/server/mocks/storage.go -package=mock Storage

lint:
	go fmt ./...
	find . -name '*.go' -exec goimports -local github.com/pisarevaa/metrics -w {} +
	find . -name '*.go' -exec golines -w {} -m 120 \;
	golangci-lint run ./...

coverage_report:
	go test -coverpkg=./... -count=1 -coverprofile=.coverage.out ./...
	go tool cover -html .coverage.out -o .coverage.html