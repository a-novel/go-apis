COVER_FILE=$(CURDIR)/coverage.out
BIN_DIR=$(CURDIR)/bin

PKG="github.com/a-novel/go-apis"

PKG_LIST=$(shell go list $(PKG)/... | grep -v /vendor/)

# Runs the test suite.
test:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --junitfile report.xml --format pkgname -- -count=1 -p 1 -v -coverpkg=./...

# Runs the test suite in race mode.
race:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		gotestsum --packages="./..." --format pkgname -- -race -count=1 -p 1 -v -coverpkg=./...

# Run the test suite in memory-sanitizing mode. This mode only works on some Linux instances, so it is only suitable
# for CI environment.
msan:
	POSTGRES_URL=$(POSTGRES_URL_TEST) ENV="test" \
		env CC=clang env CXX=clang++ gotestsum --packages="./..." --format testname -- -msan -short $(PKG_LIST) -p 1

.PHONY: all test race msan
