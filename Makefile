.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./...

.PHONY: generate-mocks
generate-mocks: ## Generate mocks by considering ./mocks/mockgen.go.
	go generate ./...

.PHONY: docker-run
docker-run: test
	docker build -t zt-event-logger . \
		&& docker run --env-file .env-docker -p 8080:8080 zt-event-logger
