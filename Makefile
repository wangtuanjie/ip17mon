all: test build

test:
	go test -v -race .

build:
	@cd example/qip && bash build.sh && mv qip ../../
