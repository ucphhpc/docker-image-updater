SHELL:=/bin/bash
PACKAGE_NAME=docker-image-updater
PACKAGE_NAME_FORMATTED=$(subst -,_,$(PACKAGE_NAME))
OWNER=ucphhpc
SERVICE_NAME=docker-image-updater
IMAGE=$(PACKAGE_NAME)
DOCKER_COMPOSE=$(shell which docker-compose || echo 'docker compose')
# Enable that the builder should use buildkit
# https://docs.docker.com/develop/develop-images/build_enhancements/
DOCKER_BUILDKIT=1
TAG=edge
ARGS=

.PHONY: all init dockerbuild dockerclean dockerpush daemon down clean dist distclean maintainer-clean
.PHONY: install uninstall installcheck check

all: init dockerbuild

init:
ifeq ($(shell test -e defaults.env && echo yes), yes)
ifneq ($(shell test -e .env && echo yes), yes)
		ln -s defaults.env .env
endif
endif

dockerbuild:
	${DOCKER_COMPOSE} build ${ARGS}

dockerclean:
	docker rmi -f $(OWNER)/$(IMAGE):$(TAG)

dockerpush:
	docker push $(OWNER)/$(IMAGE):$(TAG)

daemon:
	docker stack deploy -c <(${DOCKER_COMPOSE} config) $(SERVICE_NAME) $(ARGS)

down:
	docker stack rm $(SERVICE_NAME) $(ARGS)

clean:
	$(MAKE) dockerclean
	$(MAKE) distclean
	rm -fr .env
	rm -fr .pytest_cache
	rm -fr tests/__pycache__

dist:
### PLACEHOLDER ###

distclean:
### PLACEHOLDER ###

maintainer-clean:
	@echo 'This command is intended for maintainers to use; it'
	@echo 'deletes files that may need special tools to rebuild.'

install-dep:
### PLACEHOLDER ###

install:
	$(MAKE) install-dep

uninstall:
### PLACEHOLDER ###

uninstallcheck:
### PLACEHOLDER ###

installcheck:
### PLACEHOLDER ###

check:
### PLACEHOLDER ###