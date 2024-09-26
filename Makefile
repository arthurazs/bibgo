.PHONY: run test

run:
	go run cmd/cli/main.go

test:
	go test -v
