## Include variables from the .envrc file
include .envrc


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ / '

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/backend/... -db-dsn=${DB_DSN} -llm-api-key=${LLM_API_KEY}

## run/dataloader: run the cmd/dataloader application
.PHONY: run/dataloader
run/dataloader:
	@go run ./cmd/dataloader/... -db-dsn=${DB_DSN} -llm-api-key=${LLM_API_KEY}

## run/json-converter: run the cmd/json-converter application
.PHONY: run/json-converter
run/json-converter:
	@go run ./cmd/json-converter/... 

## run/similarity-search-example: run the cmd/similarity-search-example application
.PHONY: run/similarity-search-example
run/similarity-search-example:
	@go run ./cmd/similarity-search-example/... -db-dsn=${DB_DSN} --llm-api-key=${LLM_API_KEY} $(GO_FLAGS) 


## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	@read -p "Enter migration name: " name; \
	goose -dir ./migrations create $${name} sql 	

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	goose -dir ./migrations postgres "${DB_DSN}" up

## db/migrations/down: rollback migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Rolling back migrations'
	goose -dir ./migrations postgres "${DB_DSN}" down

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api .
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api .
