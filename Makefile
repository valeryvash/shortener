.PHONY: statictest build test check clean

statictest:
	STATIC_TEST=$$(command -v statictest); \
	echo "Using STATIC_TEST: $$STATIC_TEST"; \
	go vet -v -vettool=$$STATIC_TEST ./...

build:
	@go build -o ./cmd/shortener/shortener ./cmd/shortener/main.go

test_v2:
	SERVER_PORT=$$(random unused-port); \
	echo "Using port: $$SERVER_PORT"; \
	shortenertest_v2 -test.v -test.run="^TestIteration[1-4]$$" \
	-source-path=. \
	-binary-path=cmd/shortener/shortener \
	-server-port=$$SERVER_PORT

check: statictest build test_v2 clean

clean:
	rm -f ./cmd/shortener/shortener

generate:
	@go generate ./...
