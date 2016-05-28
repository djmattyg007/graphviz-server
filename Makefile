GO=$(shell which go)

build:
	CGO_ENABLED=0 $(GO) build graphviz-server.go
