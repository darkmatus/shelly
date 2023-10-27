#!/bin/bash
#
# Test application

set -e

source $(dirname $0)/base.sh

service_endpoint="http://localhost:8021/health"

info "Bringing services up with docker compose"
docker-compose -f docker-compose.yml up -d

# On error or when finished, shutdown containers
trap 'gracefull_shutdown_docker_compose "-f docker-compose.yml"' ERR EXIT

# check readiness before running tests
wait_for_endpoint $service_endpoint

info "Running migrations"
  docker-compose exec -T app go run . migration
info "migrations completed"

info "Running audit"
	docker-compose exec -T app bash -c "go vet ./... &&  go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./... && go mod verify"
info "audit completed"

info "Running tests"
 docker-compose exec -T app bash -c "mkdir -p reports && go test -race -coverprofile=reports/coverage.out -vet=off ./... && gocover-cobertura < reports/coverage.out > reports/cobertura-coverage.xml && go tool cover -func=reports/coverage.out"
info "tests completed"