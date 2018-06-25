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

.PHONY: docker-binary
docker-binary:
	go build -ldflags "$(LDFLAGS)" -o gitea-github-migrator

.PHONY: release
release:
	@hash gox > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	gox -ldflags "$(LDFLAGS)" -output "releases/gitea-github-migrator_{{.OS}}_{{.Arch}}"

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
	go test -cover ./...
	