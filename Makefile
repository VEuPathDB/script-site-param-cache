build:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/param-cache cmd/main.go
