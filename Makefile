build:
	GOOS=darwin go build -o client client/main.go
	GOOS=linux go build -o server server/main.go