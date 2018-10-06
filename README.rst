
====================
docker-image-updater
====================

---------------
Getting Started
---------------

Either use docker to build an image, or install go and run the updater.
Through go setup check the debian image for updates every 10 minutes (the default if left out) ::

    go build
    ./docker-image-updater -image debian -update-interval 10

Build as a docker image (defaults to use the :edge tag)::

    make build

Which produces an image called rasmunk/docker-image-updater:edge, override the TAG variable in the makefile to change this.
To run an updater container that continuously checks for updates against the debian image every 10 minutes ::

    docker run --mount type=bind,src=/var/run/docker.sock,target=/var/run/docker.sock rasmunk/docker-image-updater:edge -image debian

