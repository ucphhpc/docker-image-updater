SHELL=/bin/bash
PACKAGE_NAME=docker-image-updater
OWNER?=ucphhpc
SERVICE_NAME=${PACKAGE_NAME}
IMAGE=${PACKAGE_NAME}

# Enable that the builder should use buildkit
# https://docs.docker.com/develop/develop-images/build_enhancements/
DOCKER_BUILDKIT=1
# NOTE: dynamic lookup with docker as default and fallback to podman
DOCKER = $(shell which docker 2>/dev/null || which podman 2>/dev/null)
# if docker compose plugin is not available, try old docker-compose/podman-compose
ifeq (, $(shell ${DOCKER} help|grep compose))
	DOCKER_COMPOSE = $(shell which docker-compose 2>/dev/null || which podman-compose 2>/dev/null)
else
	DOCKER_COMPOSE = ${DOCKER} compose
endif
$(echo ${DOCKER_COMPOSE} >/dev/null)

-include .env
ARGS=


.PHONY: all
all: .env dockerbuild

.env:
	@echo
	@echo "*** No environment selected - using default ***"
	@echo
	ln -s defaults.env .env
	@sleep 2

.PHONY: build
build: dockerbuild

.PHONY: dockerbuild
dockerbuild:
	${DOCKER_COMPOSE} build ${ARGS}

.PHONY: dockerclean
dockerclean:
	docker rmi -f $(OWNER)/$(IMAGE):$(TAG)

.PHONY: dockerpush
dockerpush:
	docker push $(OWNER)/$(IMAGE):$(TAG)

.PHONY: deamon
daemon:
	docker stack deploy -c <(${DOCKER_COMPOSE} config) ${SERVICE_NAME} $(ARGS)

daemon-down:
	docker stack rm $(SERVICE_NAME)

.PHONY: up
up:
	${DOCKER_COMPOSE} up -d $(ARGS)

.PHONY: down
down:
	${DOCKER_COMPOSE} down $(ARGS)

.PHONY: clean
clean: dockerclean distclean
	rm -fr .env
	rm -fr .pytest_cache
	rm -fr tests/__pycache__

.PHONY: dist
dist:
### PLACEHOLDER ###

.PHONY: distclean
distclean:
### PLACEHOLDER ###

.PHONY: maintainer-clean
maintainer-clean:
	@echo 'This command is intended for maintainers to use; it'
	@echo 'deletes files that may need special tools to rebuild.'

.PHONY: install-dep
install-dep:
### PLACEHOLDER ###

.PHONY: install
install: install-dep
### PLACEHOLDER ###

.PHONY: uninstall
uninstall:
### PLACEHOLDER ###

.PHONY: uninstallcheck
uninstallcheck:
### PLACEHOLDER ###

.PHONY: installcheck
installcheck:
### PLACEHOLDER ###

.PHONY: check
check:
### PLACEHOLDER ###