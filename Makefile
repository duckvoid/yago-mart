BIN_NAME= gophermart
BIN_PATH = ./cmd/gophermart/
MOCK_DIR = ./mocks

VERBOSE ?= false
ifeq ($(VERBOSE), true)
  METRICSTEST_VERBOSE= -test.v
else
  METRICSTEST_VERBOSE=
endif

RUN_PORT=8080
RUN_ADDRESS="localhost"

ACCRUAL_PORT=8081
ACCRUAL_ADDRESS="localhost"
DB_STRING="postgres://postgres:postgres@localhost:5432/mart?sslmode=disable"

.PHONY: build build_gophermart mocks test audit test/cover upgradeable lint vuln tidy

build: tidy build_gophermart
	@echo "gophermart built"

build_gophermart:
	go build -v -o $(BIN_PATH)/$(BIN_NAME) $(BIN_PATH)

test: audit
	go test -v -race -buildvcs ./...

test/cover: audit
	go test -v -race -buildvcs -coverprofile=./.coverage.out ./...
	go tool cover -html=./.coverage.out

audit: upgradeable mocks lint vuln
	@echo "Audit passed"

mocks:
	@echo "Mocks generating"
	@go install github.com/gojuno/minimock/v3/cmd/minimock@latest
	@mkdir -p "./mocks"
	@go generate ./internal/...

tidy:
	@echo "Running tidy"
	go mod tidy --diff
	go mod verify

vuln:
	@echo "Running vuln check"
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

lint:
	@echo "Running golangci-lint"
	golangci-lint run ./...

upgradeable:
	@go list -u -f \
		'{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' \
		-m all; echo "This packages need to upgrade" || echo "No need to upgrade"

clean:
	@rm -rf ${MOCK_DIR} .coverage.out ./internal/scraper/testdata

gophermarttest: build
	@echo "All metric tests passed"
	gophermarttest \
		-test.v -test.run=^TestGophermart$ \
		-gophermart-binary-path=$(BIN_PATH)/$(BIN_NAME) \
		-gophermart-host=$(RUN_ADDRESS) \
		-gophermart-port=$(RUN_PORT) \
		-gophermart-database-uri=$(DB_STRING) \
		-accrual-binary-path=cmd/accrual/accrual_darwin_amd64 \
		-accrual-host=$(ACCRUAL_ADDRESS) \
		-accrual-port=$(ACCRUAL_PORT) \
		-accrual-database-uri=$(DB_STRING)