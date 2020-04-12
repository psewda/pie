VERSION=1.0
BUILD=1
BASE_DIR=$(shell pwd)
OUTPUT=$(BASE_DIR)/bin
APP=pie
MAIN=$(BASE_DIR)/cmd/$(APP)/main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	$(eval LDFLAGS="-s -w $(shell $(BASE_DIR)/ci/version.sh 'linux/amd64')")
	go build -ldflags=$(LDFLAGS) -o $(OUTPUT)/linux-amd64/$(APP) $(MAIN)
	
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64
	$(eval LDFLAGS="-s -w $(shell $(BASE_DIR)/ci/version.sh 'windows/amd64')")
	go build -ldflags=$(LDFLAGS) -o $(OUTPUT)/windows-amd64/$(APP).exe $(MAIN)

build: build-linux build-windows

test:
	go test -v ./...

run:
	go build -o $(OUTPUT)/$(APP) $(MAIN)
	./bin/pie
