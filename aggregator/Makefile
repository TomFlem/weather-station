.PHONY: build

GO=CGO_ENABLED=0 GO111MODULE=on go
GOFLAGS=

AGGREGATOR=weather-aggregator

.PHONY: $(EDGEBUILDERMETRICS)

build: $(EDGEBUILDERMETRICS)
	go vet ./...
	$(GO) build $(GOFLAGS) -o weather-aggregator ./cmd/aggregator/main.go

vet: $(EDGEBUILDERMETRICS)
	go mod tidy
	go vet ./...