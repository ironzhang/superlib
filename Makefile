# Makefile

TestPackages=$(shell go list ./...)

.PHONY: test

all: test

test:
	go test $(TestPackages)

