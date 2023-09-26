run:
	go run main.go

dev:
	nodemon --exec go run main.go

install:
	go mod tidy