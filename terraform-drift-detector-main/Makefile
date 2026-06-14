.PHONY: build test run-server run-scan clean

export PATH := /usr/local/go/bin:$(PATH)

build:
	go build -o bin/tfdriftctl ./cmd/tfdriftctl
	go build -o bin/drift-server ./cmd/drift-server

test:
	go test ./...

run-server: build
	./bin/drift-server -config configs/tfdriftctl.yaml

run-scan: build
	./bin/tfdriftctl scan --state testdata/state/sample.tfstate --provider aws --skip-cloud

clean:
	rm -rf bin/
