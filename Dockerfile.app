## Dockerfile for a Node.js application
##
## == How to use
## Build the Docker image using the following command: `docker build --no-cache  -t local/source2adoc:dev -f Dockerfile.app .`
##
## Use `docker run --rm -p "3000:9000" sommerfeldio/source2adoc:rc` to run the most
## recent release candidate from DockerHub.
##
## @see docker-compose.yml


## This stage is based on the official Node.js image and copies the application code into the
## image, installs the dependencies and bundles the application
FROM node:22.3.0-alpine3.20 AS build
LABEL maintainer="sebastian@sommerfeld.io"


## Expose the documentation site using httpd.
##
## To avoid running the httpd and thes image as `root`, the permissions of `/usr/local/apache2/logs`
## are changed to allow `www-data` to write logs. Additionally the default http port is changed to
## `3000`, so keep that in mind when mapping ports in a `docker run ...` command. This way the image
## can be used without root permissions because the httpd server inside the container is started
## with the already existing user `www-data`.
##
## The webserver exposes his status information through http://localhost:3000/server-status.
FROM httpd:2.4.59-alpine3.19 AS run
LABEL maintainer="sebastian@sommerfeld.io"

ARG USER=www-data
RUN chown -hR "$USER:$USER" /usr/local/apache2 \
    && chmod g-w /usr/local/apache2/conf/httpd.conf \
    && chmod g-r /etc/shadow \
    && rm /usr/local/apache2/htdocs/index.html

# TODO ... A CLI app should not run in a webserver !!!!!
RUN echo "A CLI app should not run in a webserver !!!!!" > /usr/local/apache2/htdocs/index.html

USER "$USER"
EXPOSE 3000

