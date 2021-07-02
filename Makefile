GO_PACKAGES = $(shell go list ./... | grep -v vendor)
GO_FILES = $(shell find . -name "*.go" | grep -v vendor | uniq)
LOCAL_PACKAGES="github.com/"
MOCKERY = ${GOPATH} mockery

format:
	@echo "==> Formatting"
	@gofmt -s -l -w $(GO_FILES)
	@goimports -w -local $(LOCAL_PACKAGES) $(GO_FILES)

lint:
	golangci-lint run ${GOLANGCI_CONCURRENCY} --sort-results --out-format colored-line-number ${LINT_FLAGS} ${LINT_TARGET}

vet:
	@echo "==> Vetting for linux"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go vet ${GO_PACKAGES}

cover:
	@echo "==> Testing coverage"
	@goverage -v -covermode count -coverprofile=cover.out $(GO_PACKAGES) 2>&1 | grep -v "warning: no packages "
	@go tool cover -html=cover.out

# run all benchmarks but no tests
bench:
	@echo "==> Benchmarking"
	@go test -run=XXX -bench=. $(GO_PACKAGES)

# install supports installing binaries to the $GOPATH/bin for the host platform
install:
	@echo "==> Installing cmds to $$GOPATH/bin"
	@go install ${LDFLAGS} ./cmd/...

# Installs and runs manager locally
start: install
	$(GOPATH)/bin/booking

start-local:
	HOST=localhost PORT=8081 MAXTIMEBLOCK=60 LOGLEVEL=trace DBURL='${DATABASE_URL} SWAGGERDIST=/home/jpetersen/Workspace/swagger-ui/dist go run cmd/booking/main.go

docker-build:
	docker build -t booking .

docker-start:
	docker run -e HOST=127.0.0.1 -e PORT=8081 -e LOGLEVEL=trace -e DBURL=${DATABASE_URL} -e SWAGGERDIST=/swagger-ui/dist -p 8081:8081 -it --name booking-service booking

docker-start-local:
	docker run -e HOST=127.0.0.1 -e PORT=8081 -e LOGLEVEL=trace -e DBURL=${DATABASE_URL} -e SWAGGERDIST=/swagger-ui/dist -p 8081:8081 -it --name booking-service booking

mocks:
	@echo "==> Generating mocks"
	@${MOCKERY} --dir=./repository --name=Repository --output=./repository/mocks
	@${MOCKERY} --dir=./service --name=RoomService --output=./service/mocks
	@${MOCKERY} --dir=./service --name=BookingService --output=./service/mocks

test:
	go test -v -coverprofile=coverage.out -timeout=1m -race ./...