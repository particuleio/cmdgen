BINARY_NAME=cmdgen
BINARY_DIR=bin
RUN_ARGS?=start ./files/test.yaml 

define build_project
	mkdir -p ${BINARY_DIR}
	go build -o ${BINARY_DIR} ./...
endef

# https://dev.to/victoria/how-to-create-a-self-documenting-makefile-2b0e
.DEFAULT_GOAL := help 
help: # Show this help
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'


deps: # Install dependencies
	go get github.com/fatih/color
	go get gopkg.in/yaml.v3

build: # Build project
	$(call build_project)

test: # Test project
	@go test -timeout 30s github.com/particuleio/cmdgen/pkg/cmdgen

run: # Build and run the project
	$(call build_project)
	./${BINARY_DIR}/${BINARY_NAME} ${RUN_ARGS}

clean: # Removes object files from package source directories and removes build dir
	go clean
	rm -rf ${BINARY_DIR}
