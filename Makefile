SHELL := /bin/sh

dev:
	docker-compose -f internal/kit/docker/docker-compose.yml up
