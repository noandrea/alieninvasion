GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)
GIT_DESCR = $(shell git describe --tags --always) 
APP=alieninvasion
# build output folder
OUTPUTFOLDER = dist
# docker image
DOCKER_REGISTRY = gitlab
DOCKER_IMAGE = alieninvasion
DOCKER_TAG = $(shell git describe --always --tags)
# build paramters
OS = linux
ARCH = amd64

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

default: build

workdir:
	mkdir -p dist

build: build-dist

build-dist: $(GOFILES)
	@echo build binary to $(OUTPUTFOLDER)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "-w -s -X main.Version=$(GIT_DESCR)" -o $(OUTPUTFOLDER)/$(APP) .
	@echo copy resources
	cp -r README.md LICENSE $(OUTPUTFOLDER)
	@echo done

build-relase:
	@echo build binary to $(OUTPUTFOLDER)
	$(eval OS=darwin)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "-X main.Version=$(GIT_DESCR)" -o $(OUTPUTFOLDER)/$(OS)/$(APP) .
	$(eval OS=windows)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "-X main.Version=$(GIT_DESCR)" -o $(OUTPUTFOLDER)/$(OS)/$(APP).exe .
	$(eval OS=linux)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "-X main.Version=$(GIT_DESCR)" -o $(OUTPUTFOLDER)/$(OS)/$(APP) .
	@echo copy resources
	cp -r README.md LICENSE CHANGELOG.md $(OUTPUTFOLDER)
	@echo done

test: test-all

test-all:
	@go test -v $(GOPACKAGES) -coverprofile .testCoverage.txt

bench: bench-all

bench-all:
	@go test -bench -v $(GOPACKAGES)

lint: lint-all

lint-all:
	@golint -set_exit_status $(GOPACKAGES)

clean:
	@echo remove $(OUTPUTFOLDER) folder
	@rm -rf $(OUTPUTFOLDER)
	@echo done

docker: docker-build

docker-build:
	@echo copy resources
	docker build -t $(DOCKER_IMAGE)  .
	@echo done

docker-push:
	@echo push image
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	aws ecr get-login --no-include-email --region eu-central-1 --profile aeternity-sdk | sh
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo done

docker-run: 
	docker run -p 1905:1905 $(DOCKER_IMAGE):latest

debug-start:
	@go run main.go start

changelog:
	git-chglog --output CHANGELOG.md

git-release:
	@echo making release
	git tag $(GIT_DESCR)
	git-chglog --output CHANGELOG.md
	git tag $(GIT_DESCR) --delete
	git add CHANGELOG.md && git commit -m "v$(GIT_DESCR)" -m "Changelog: https://github.com/noandrea/alieninvasion/blob/master/CHANGELOG.md"
	git tag $(GIT_DESCR)
	git push --tags
	@echo release complete


_release-patch:
	$(eval GIT_DESCR = $(shell git describe --tags | awk -F '("|")' '{ print($$1)}' | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g'))
release-patch: _release-patch git-release

_release-minor:
	$(eval GIT_DESCR = $(shell git describe --tags | awk -F '("|")' '{ print($$1)}' | awk -F. '{$$(NF-1) = $$(NF-1) + 1;} 1' | sed 's/ /./g' | awk -F. '{$$(NF) = 0;} 1' | sed 's/ /./g'))
release-minor: _release-minor git-release

_release-major:
	$(eval GIT_DESCR = $(shell git describe --tags | awk -F '("|")' '{ print($$1)}' | awk -F. '{$$(NF-2) = $$(NF-2) + 1;} 1' | sed 's/ /./g' | awk -F. '{$$(NF-1) = 0;} 1' | sed 's/ /./g' | awk -F. '{$$(NF) = 0;} 1' | sed 's/ /./g' ))
release-major: _release-major git-release 

