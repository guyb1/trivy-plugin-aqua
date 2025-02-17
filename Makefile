SHELL=/usr/bin/env bash


.PHONY: test
test:
	if cat ./plugin.yaml | grep uri: | cut -c 10-200 | xargs -n 1 curl -o /dev/null --silent --head --write-out '%{http_code}\n'  -L | grep -q "404"; then echo 'Artifacts links are broken'; exit 1;else echo 'Artifacts links are valid';fi
	go test -v ./...

.PHONY: integration-test
integration-test:
	curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.27.0
	mkdir -p /home/runner/.trivy/plugins/aqua/
	go build -o /home/runner/.trivy/plugins/aqua/aqua cmd/aqua/main.go
	cp plugin.yaml /home/runner/.trivy/plugins/aqua/
	trivy config .
	trivy fs --debug --security-checks config,vuln,secret .
	docker pull alpine
	trivy --debug image alpine

.PHONY: update-plugin
update-plugin:
	@LEVEL=$(LEVEL) ./scripts/update_plugin.sh

.PHONY: proto
proto:
	pushd pkg/proto && protoc --twirp_out=. --go_out=. ./buildsecurity.proto

.PHONY: build-image
build-image:
	docker run \
  --rm \
  -e GOARCH=amd64 \
  -e GOOS=linux \
  -w /build \
  -v `pwd`:/build \
  golang:1.18 \
  go build -o bin/aqua cmd/aqua/main.go|| exit 1

.PHONY: test-image
test-image:
	docker run \
  --rm \
  -e GOARCH=amd64 \
  -e GOOS=linux \
  -w /build \
  -v `pwd`:/build \
  golang:1.18 \
  go test -v ./...|| exit 1
