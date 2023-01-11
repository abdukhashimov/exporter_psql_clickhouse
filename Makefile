include .build_info

CURRENT_DIR=$(shell pwd)

TAG=latest
ENV_TAG=latest

APP_CMD_DIR=${CURRENT_DIR}/cmd

.PHONY: build
build:
	@echo ⌛ building...
	go build -o bin/exporter cmd/main.go
	@echo ✅ building done

.PHONY: run
run:
	@echo initializing...
	go run cmd/main.go

.PHONY: clean
clean:
	@echo ⌛ cleaning...
	rm -rf bin/*
	@echo ✅ cleaning done

.PHONY: build-image
build-image:
	@echo ⌛ building the docker image...
	docker build --build-arg APP=${APP} --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	@echo ✅ building the docker image done

.PHONY: push-image
push-image:
	@echo ⌛ pushing the docker image...
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	@echo ✅ pushing the docker image done

.PHONY: service-up
service-up:
	@echo ⌛ service up...
	