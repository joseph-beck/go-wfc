.PHONY: run

run:
	scripts/run.sh

.PHONY: test

test:
	go clean -testcache 
	go mod tidy
	go test -cover ./...