APP = amamonitor
COMMANDS = fetcher server
BASE_DIR = github.com/oinume/amamonitor
VENDOR_DIR = vendor
PROTO_GEN_DIR = proto-gen
GO_GET ?= go get
GO_TEST ?= go test -v -race
GO_TEST_PACKAGES = $(shell go list ./... | grep -v vendor | grep -v e2e)
GOPATH = $(shell go env GOPATH)
LINT_PACKAGES = $(shell go list ./...)
IMAGE_TAG ?= latest
VERSION_HASH_VALUE = $(shell git rev-parse HEAD | cut -c-7)
PID = $(APP).pid


all: build

.PHONY: setup
setup: install-lint install-tools

.PHONY: install-tools
install-tools:
	cd tools && go install \
		github.com/xo/xo \
		github.com/pressly/goose/cmd/goose
#	GO111MODULE=off $(GO_GET) bitbucket.org/liamstask/goose/cmd/goose

.PHONY: tools
tools: ## install dependent tools
	cd tools && go install \
		github.com/kouzoh/wrench \
		go.mercari.io/yo

.PHONY: install-lint
install-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.17.1

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: install
install:
	go install github.com/oinume/lekcije/server/cmd/lekcije

.PHONY: git-config
git-config:
	echo "" > ~/.gitconfig
	git config --global url."https://github.com".insteadOf git://github.com
	git config --global http.https://gopkg.in.followRedirects true

.PHONY: build
build: $(foreach command,$(COMMANDS),build/$(command))

# TODO: find server/cmd -type d | xargs basename
# OR CLIENTS=hoge fuga proto: $(foreach var,$(CLIENTS),proto/$(var))
build/%:
	CGO_ENABLED=0 GO111MODULE=on go build -o bin/$* $(BASE_DIR)/backend/cmd/$*

clean:
	${RM} $(foreach command,$(COMMANDS),bin/$(command))

.PHONY: ngrok
ngrok:
	ngrok http -subdomain=amamonitor -host-header=localhost 4000

.PHONY: test
test: go/test

.PHONY: go/test
go/test:
	$(GO_TEST) $(GO_TEST_PACKAGES)

.PHONY: goimports
goimports:
	goimports -w ./server ./e2e

.PHONY: lint
lint: go/lint

.PHONY: go/lint
go/lint:
	golangci-lint run

.PHONY: docker/build
docker/build: $(foreach command,$(COMMANDS),docker/build/$(command))

.PHONY: docker/build/%
docker/build/%:
	docker build --pull --no-cache \
	-f docker/Dockerfile-$* \
	--tag asia.gcr.io/amamonitor/$*:$(IMAGE_TAG) .

.PHONY: gcloud/builds
gcloud/builds: $(foreach command,$(COMMANDS),gcloud/builds/$(command))

.PHONY: gcloud/builds/%
gcloud/builds/%:
	gcloud builds submit . \
	--project $(GCP_PROJECT_ID) \
	--config=gcloud-builds.yml \
	--substitutions=_IMAGE_TAG=$(IMAGE_TAG),_COMMAND=$*

.PHONY: db/goose/%
db/goose/%:
	goose mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/amamonitor?charset=utf8mb4&parseTime=true&loc=UTC" $*

.PHONY: reset-db
reset-db:
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS amamonitor"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot -e "DROP DATABASE IF EXISTS amamonitor_test"
	mysql -h $(DB_HOST) -P 13306 -uroot -proot < db/create_database.sql

kill:
	kill `cat $(PID)` 2> /dev/null || true

restart: kill clean build/server
	bin/server & echo $$! > $(PID)

watch: restart
	fswatch -o -e ".*" -e vendor -e node_modules -e .venv -i "\\.go$$" . | xargs -n1 -I{} make restart || make kill
