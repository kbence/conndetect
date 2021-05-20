BIN = out/conndetect
ARCH =

SOURCES = $(call rwildcard,./,*.go)
PACKAGES_WITH_TESTS = $(shell find . -name '*_test.go' | sed 's/[^/]\+\.go//' | sort | uniq)

.PHONY: all
all: $(BIN)

.PHONY: test
test:
	./test.sh

$(BIN): $(SOURCES)
	echo $^
	mkdir -p ./out
	GOARCH=$(ARCH) go build -o "$@" .

# Functions
rwildcard=$(wildcard $1$2) $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2))
