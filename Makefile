.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create-user functions/users/create-user.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/login functions/sessions/login.go
clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
