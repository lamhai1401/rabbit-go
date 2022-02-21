SHELL := /bin/bash

# 0
run_server:
	docker run --name test-rabbit -e RABBITMQ_DEFAULT_USER=username -e RABBITMQ_DEFAULT_PASS=password rabbitmq:3-management \
	go run .
.PHONY: run_server