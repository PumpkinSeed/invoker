GOMODULES ?= ./...
GOFILES ?= $(shell go list $(GOMODULES) | grep -v /vendor/)

format:
	@echo "--> Running go fmt"
	@go fmt $(GOFILES)