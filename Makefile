build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/manager manager/main.go
