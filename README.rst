====================
docker-image-updater
====================

Docker image updater that continuously checks for updates for a given docker image,
has two options for execution:

- build and execute the binary native on the host itself.
- run as a container with the hosts docker socket mounted inside the container.

---------------
Getting Started
---------------

Either use docker to build an image, or build the binary and run it on the host.
Through go setup check the debian image for updates every 10 minutes (the default if left out)::

    go build ./...
    ./docker-image-updater -update debian -interval 10
    
    # If only updated images should be kept
    ./docker-image-updater -update debian -prune
    
    # If updated untagged images should be pruned aswell
    ./docker-image-updater -update debian -prune -prune-untagged
    
    # If in addition to images being updated other images should just be kept without being updated
    ./docker-image-updater -update debian -protect ubuntu -prune

Build as a docker image (defaults to use the :edge tag)::

    make build
    
    # override the build tag, e.g
    make build TAG=latest

Which produces an image called ucphhpc/docker-image-updater:edge by default,
override the TAG variable in the makefile to change this. To run an updater container
that continuously checks for updates against the debian image every 10 minutes::

    docker run --mount type=bind,src=/var/run/docker.sock,target=/var/run/docker.sock ucphhpc/docker-image-updater:edge -update debian


