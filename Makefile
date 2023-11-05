mod:
	go mod tidy
	go mod vendor

test:
	go test -mod=vendor ./... -v