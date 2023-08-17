# Go parameters
BUILD_DIR=bin
GOCMD=GOPRIVATE=$(GOPRIVATE) go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test ./...
GOGET=$(GOCMD) get
PROTOCMD=proto
GOLINT=golangci-lint run
VERSION=$$(eval "git describe 2> /dev/null || echo v0.0.0-pre-$$(git rev-parse --short HEAD)")
GOVENDOR=test -d vendor || $(GOCMD) mod vendor # It is unlikely to get the private go packages in a container, so it is advised to run `go mod vendor` before running `docker build`
GOBUILD=$(GOCMD) build -o ${BUILD_DIR} -ldflags "-X main.version=${VERSION}"

all: main docker

.PHONY: lint
lint:
	${GOLINT}

.PHONY: test
test:
	${GOTEST} -v

.PHONY: coverage
coverage:
	$(GOCMD) test -coverprofile=coverage.out -covermode=count ./...
	$(GOCMD) tool cover -func=coverage.out
	$(GOCMD) tool cover -html=coverage.out

.PHONY: main
main:
	@echo "Building main"
	@mkdir -p ${BUILD_DIR}
	${GOVENDOR}
	${GOBUILD}

.PHONY: docker
docker:
	@echo "Building docker container for backend component"
	docker build . -t k8s_sample_backend
