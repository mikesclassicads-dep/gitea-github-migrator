IMPORT := git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator
GO ?= go


ifneq ($(DRONE_TAG),)
	VERSION ?= $(subst v,,$(DRONE_TAG))
else
	ifneq ($(DRONE_BRANCH),)
		VERSION ?= $(subst release/v,,$(DRONE_BRANCH))
	else
		VERSION ?= master
	endif
endif

LDFLAGS := -X main.version=$(VERSION) -X main.build=$(DRONE_BUILD_NUMBER)

.PHONY: all
all:

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o gitea-github-migrator

.PHONY: build-binary-web
build-binary-web:
	go build -ldflags "$(LDFLAGS)" -tags web -o gitea-github-migrator

.PHONY: build-web
build-web: packr build-binary-web packr-clean

.PHONY: packr
packr:
	@hash packr > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/gobuffalo/packr/...; \
	fi
	packr -z

.PHONY: packr-clean
packr-clean:
	@hash packr > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/gobuffalo/packr/...; \
	fi
	packr clean

.PHONY: clean
clean: packr-clean
	go clean ./...

.PHONY: docker-binary
docker-binary: build

.PHONY: docker-binary-web
docker-binary-web: build-web


.PHONY: generate-release-file
generate-release-file:
	echo $(VERSION) > .version

.PHONY: release
release: packr release-builds packr-clean

.PHONY: release-builds
release-builds:
	@hash gox > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	gox -ldflags "$(LDFLAGS)" -tags web -output "releases/gitea-github-migrator_{{.OS}}_{{.Arch}}"

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	golint -set_exit_status $(go list ./...)

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: lint vet
	go test -tags web -cover ./...
	