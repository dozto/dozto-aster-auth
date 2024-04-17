build:
	go build -o ./out/gpt
dev:
	GO_ENV=dev air
tidy:
	go mod tidy -v
test:
	GO_ENV=test go test ./app/... ./pkg/... -cover
test-v:
	GO_ENV=test go test -v ./app/... ./pkg/...
test-clean:
	go clean -testcache