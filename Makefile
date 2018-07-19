APPLICATION_VERSION := "0.1.0"
COMMIT_HASH := `git log -1 --format="%H"`
BUILD_TIMESTAMP := `date +%s`

DEP := $(shell command -v dep 2> /dev/null)
GIT := $(shell command -v git 2> /dev/null)
GO := $(shell command -v go 2> /dev/null)
YARN := $(shell command -v yarn 2> /dev/null)
ZIP := $(shell command -v zip 2> /dev/null)

LDFLAGS :=-ldflags "-X github.com/dotStart/Watchdog/metadata.version=${APPLICATION_VERSION} -X github.com/dotStart/Watchdog/metadata.commitHash=${COMMIT_HASH} -X \"github.com/dotStart/Watchdog/metadata.buildTimstampRaw=${BUILD_TIMESTAMP}\""
PLATFORMS := darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386/.exe windows/amd64/.exe

# magical formula:
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

all: print-config check-env install-deps build-ui generate-sources $(PLATFORMS)

print-config:
	@echo "==> Build Configuration"
	@echo ""
	@echo "  Application Version: $(APPLICATION_VERSION)"
	@echo "          Commit Hash: $(COMMIT_HASH)"
	@echo "      Build Timestamp: $(BUILD_TIMESTAMP)"
	@echo ""

check-env:
	@echo "==> Verifying build environment"
	@echo ""
	@echo -n "Checking for dep ... "
ifndef DEP
	@echo "Not Found"
	$(error "dep is unavailable")
endif
	@echo $(DEP)
	@echo -n "Checking for git ... "
ifndef GIT
	@echo "Not Found"
	$(error "git is unavailable")
endif
	@echo $(GIT)
	@echo -n "Checking for go ... "
ifndef GO
	@echo "Not Found"
	$(error "go is unavailable")
endif
	@echo $(GO)
	@echo -n "Checking for yarn ... "
ifndef YARN
	@echo "Not Found"
	$(error "yarn is unavailable")
endif
	@echo $(YARN)
	@echo -n "Checking for zip ... "
ifndef ZIP
	@echo "Not Found"
else
	@echo $(ZIP)
endif
	@echo ""

.ONESHELL:
install-deps:
	@echo "==> Installing dependencies"
	@echo ""
	@$(GO) get -u github.com/jteeuwen/go-bindata/go-bindata
	@$(GO) get -u github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs
	@$(DEP) ensure -v
	@cd ui
	@$(YARN)
	@echo ""

.ONESHELL:
build-ui:
	@echo "==> Building UI"
	@cd ui
	@mkdir resources/3rdParty
	@cp node_modules/semantic-ui-css/semantic.min.css resources/3rdParty/semantic.min.css
	@echo ""

generate-sources:
	@echo "==> Generating Sources"
	@$(GO) generate -v github.com/dotStart/Watchdog/ui
	@echo ""

$(PLATFORMS):
	@echo "==> Building for ${os} (${arch})"
	@export GOOS=$(os); export GOARCH=$(arch); $(GO) build -v ${LDFLAGS} -o build/$(os)-$(arch)/watchdog$(ext)
	@echo ""
ifdef ZIP
	@$(ZIP) build/watchdog-$(APPLICATION_VERSION)-$(os)-$(arch).zip build($(os)-$(arch)/watchdog$(ext)
	@echo ""
endif

.PHONY: all
