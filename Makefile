.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/create-user functions/create-user/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login functions/login-user/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/chat functions/chat/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/read functions/read/main.go


clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
