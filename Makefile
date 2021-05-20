BIN = out/conndetect
ARCH =

DOCKER_IMAGE=conndetect

SOURCES = $(call rwildcard,./,*.go)
PACKAGES_WITH_TESTS = $(shell find . -name '*_test.go' | sed 's/[^/]\+\.go//' | sort | uniq)

.PHONY: all
all: $(BIN)

.PHONY: test
test:
	./test.sh

.PHONY: docker-image
docker-image:
	# TODO: create a repo and reference that here
	docker build -t "$(DOCKER_IMAGE)" .

.PHONY: docker-run
docker-run: docker-image
	docker run -it --rm --network host "$(DOCKER_IMAGE)"

$(BIN): $(SOURCES)
	echo $^
	mkdir -p ./out
	GOARCH=$(ARCH) go build -o "$@" .

# Functions
rwildcard=$(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))
