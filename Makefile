all:
	go clean
	go env -w GOOS=linux
	go build -a .
	go env -w GOOS=windows