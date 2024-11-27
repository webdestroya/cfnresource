gopkgs := $(shell go list ./...)

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: test
test:
	go test -v -timeout 90s $(gopkgs)

.PHONY: bench
bench:
	go test -v -bench=. $(gopkgs)

.PHONY: coverage
coverage:
	@mkdir -p coverage
	gotest -v $(gopkgs) -coverpkg=./... -coverprofile=coverage/c.out -covermode=count -short
	@cat coverage/c.out | grep -v /internal/testutils/ > coverage/c_notest.out
	@go tool cover -html=coverage/c_notest.out -o coverage/index.html
