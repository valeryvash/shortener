.PHONY: statictest build test1 check clean

statictest:
	@go vet -vettool=$$(command -v statictest) ./...

build:
	@go build -o ./cmd/shortener/shortener ./cmd/shortener/main.go

test1:
	@shortenertest -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener

test2:
	@shortenertest -test.v -test.run=^TestIteration2$ -source-path=.
	
check: statictest build test1

clean:
	rm -f ./cmd/shortener/shortener