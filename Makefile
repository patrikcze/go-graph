# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=go-api-mysql
IMAGE_NAME=go-api-mysql
CONTAINER_NAME=go-api-mysql
REGISTRY=docker.io

# MacOS or Linux
ifeq ($(shell uname -s),Darwin)
    CURL_CMD = curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-`uname -s`-`aarch64`"
else
    CURL_CMD = curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-`uname -s`-`uname -m`"
endif

# Dependency management
deps:
	$(GOGET) -u honnef.co/go/tools/cmd/staticcheck
	$(GOGET) -u github.com/mgechev/revive
	# $(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	sudo $(CURL_CMD) -o /usr/local/bin/docker-compose
	sudo chmod +x /usr/local/bin/docker-compose

# Docker parameters
IMAGE_NAME=go-api-mysql
CONTAINER_NAME=go-api-mysql

# Build the binary
build-binary: deps
	$(GOBUILD) -o $(BINARY_NAME) -v

# Build and run the container
build-docker-image: build-binary
	docker-compose build -t $(IMAGE_NAME) .
	docker-compose run -d -p 80:8080 --name $(CONTAINER_NAME) $(IMAGE_NAME)

# Stop and remove the container
stop:
	docker-compose down $(CONTAINER_NAME)
	docker-compose rm $(CONTAINER_NAME)

# Run go vet
vet:
	$(GOCMD) vet ./...

# Run golangci-lint
lint:
	golangci-lint run

# Run revive
revive:
	revive -config .revive.toml -formatter friendly ./...

# Run all checks
check: vet lint revive

# Cleanup binary file
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Build the binary and create the Docker image
build: deps
	$(GOBUILD) -o $(BINARY_NAME) -v
	docker-compose build -t $(IMAGE_NAME) .

# Push the image to the Docker registry
push:
	docker tag $(IMAGE_NAME) $(REGISTRY)/$(IMAGE_NAME)
	docker push $(REGISTRY)/$(IMAGE_NAME)

# Enable lint later
# .PHONY: build run stop vet lint revive check clean

.PHONY: build-docker-image vet revive check clean