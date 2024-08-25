TIMESTAMP=$(shell date +"%F %T %Z")
GIT_SHA=$(shell git rev-parse HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

create_version:
	printf "{\n  \"timestamp\":\"${TIMESTAMP}\",\n  \"gitSha\":\"${GIT_SHA}\",\n  \"gitBranch\":\"${GIT_BRANCH}\"\n}\n" > version.json

go_build:
	go build .

go_test:
	go test . -v

build: create_version go_build
