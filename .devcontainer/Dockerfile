## This Dockerfile is used to build the development container for the source2adoc project.
## It contains the necessary instructions to create an environment with all the dependencies
## required for development.
##
## == Prerequisites
## Having Visual Studio Code (VSCode) and the Dev Container plugin installed are
## essential prerequisites for this development environment. Docker is also mandatory.

FROM mcr.microsoft.com/devcontainers/go:1.23-bullseye
LABEL maintainer="sebastian@sommerfeld.io"

ARG USER=vscode

RUN apt-get update \
    && apt-get install -y --no-install-recommends python3-pip=20.3.4-4+deb11u1 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && pip install --no-cache-dir pre-commit==3.3.3

COPY configs/etc/login.defs /etc/login.defs
RUN chmod og-r /etc/shadow \
    && chmod 0444 /etc/login.defs

USER $USER

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN echo "alias ls='ls -a --color=auto'" >> /home/$USER/.bashrc \
    && echo "alias ll='ls -alF --color=auto'" >> /home/$USER/.bashrc
