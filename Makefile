# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	docker-compose run app go fmt ./...
	docker-compose run app go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	docker-compose run app go vet ./...
	docker-compose run app go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	docker-compose run app go mod verify

## test: run the tests
.PHONY: test
test:
	docker-compose run app go test -race -coverprofile=coverage.out -vet=off ./...
