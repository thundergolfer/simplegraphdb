 note: call scripts from /scripts
BINARY = simplegraphdb
VET_REPORT = vet.report
TEST_PORT = tests.xml
GOARCH = amd64


VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

GITHUB_USERNAME=thundergolfer
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

all: link clean test vet linux darwin windows

link:
				BUILD_DIR=${BUILD_DIR}; \
				BUILD_DIR_LINK=${BUILD_DIR_LINK}; \
				CURRENT_DIR=${CURRENT_DIR}; \
				if [ "$${BUILD_DIR_LINK}" != "$${CURRENT_DIR}" ]; then \
					echo "Fixing symlinks for build"; \
					rm -f $${BUILD_DIR}; \
					ln -s $${CURRENT_DIR} $${BUILD_DIR}; \
				fi

linux:
				cd ${BUILD_DIR}; \
				GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} . ; \
				cd - >/dev/null


darwin:
				cd ${BUILD_DIR}; \
				GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \
				cd - >/dev/null

windows:
				cd ${BUILD_DIR}; \
				GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \
				cd - >/dev/null


test:
				go test
				cd - >/dev/null

vet:
				-cd ${BUILD_DIR}; \
				godep go vet ./... > ${VET_REPORT} 2>&1 ; \
				cd - >/dev/null

fmt:
				cd ${BUILD_DIR}; \
				go fmt $${go list ./... | grep -v /vendor/) ; \
				cd - >/dev/null

clean:
				-rm -f ${TEST_REPORT}
				-rm -f ${VET_REPORT}
				-rm -f ${BINARY}-*

.PHONY: link linux darwin windows test vet fmt clean
