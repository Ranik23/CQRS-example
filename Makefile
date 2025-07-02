# Goose setup
GOOSE_VERSION := v3.8.0
BIN_DIR := ./bin
GOOSE_BIN := $(BIN_DIR)/goose
MIGRATIONS_DIR := ./migrations
GOOSE_PACKAGE := github.com/pressly/goose/v3/cmd/goose
# Mockgen setup
MOCKGEN_VERSION := v1.6.0
MOCKGEN_BIN := $(BIN_DIR)/mockgen
MOCKGEN_PACKAGE := github.com/golang/mock/mockgen

.PHONY: setup
setup: goose-setup mockgen-setup


.PHONY: docker
docker:
	go mod tidy 
	go mod vendor 
	docker compose up --build

.PHONY: k6
k6:
	k6 run test/k6/load.js

.PHONY: goose-setup
goose-setup:
	@echo "Installing goose $(GOOSE_VERSION)..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(abspath $(BIN_DIR)) go install $(GOOSE_PACKAGE)@latest
	@chmod +x $(GOOSE_BIN)
	@echo "Goose installed at $(GOOSE_BIN)"

.PHONY: swagger-setup
swagger-setup:
	@echo "Installing swagger..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(abspath $(BIN_DIR)) go install github.com/swaggo/swag/cmd/swag@latest
	@chmod +x $(BIN_DIR)/swag
	@echo "Swagger installed at $(BIN_DIR)/swag"

.PHONY: mockgen-setup
mockgen-setup:
	@echo "Installing mockgen $(MOCKGEN_VERSION)..."
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(abspath $(BIN_DIR)) go install $(MOCKGEN_PACKAGE)@latest
	@chmod +x $(MOCKGEN_BIN)
	@echo "mockgen installed at $(MOCKGEN_BIN)"


.PHONY: migrate-create
migrate-create:
	@read -p "Migration name: " name; \
	$(GOOSE_BIN) create $$name sql -dir $(MIGRATIONS_DIR)


