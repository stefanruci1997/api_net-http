build:
	go build -o bin/api_net-http
run:  build
	./bin/api_net-http
test :
	go test -v . /....