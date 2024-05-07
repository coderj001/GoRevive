# Build Stage
FROM golang:1.21 AS build-stage

LABEL app="build-GoRevive"
LABEL REPO="https://github.com/coderj001/GoRevive"

ENV PROJPATH=/go/src/github.com/coderj001/GoRevive

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/coderj001/GoRevive
WORKDIR /go/src/github.com/coderj001/GoRevive

RUN make build-alpine

# Final Stage
# FROM alpine
FROM ubuntu

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/coderj001/GoRevive"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/GoRevive/bin

WORKDIR /opt/GoRevive/bin

COPY --from=build-stage /go/src/github.com/coderj001/GoRevive/bin/gorevive /opt/GoRevive/bin/
RUN chmod +x /opt/GoRevive/bin/gorevive

# install dumb-init
# RUN apk add --no-cache dumb-init

RUN apt-get update && apt-get upgrade
RUN apt-get install -y tmux vim

ENV export EDITOR=vim
ENV export EDITOR=bash

# Create appuser
# RUN adduser -D -g '' GoRevive
# USER GoRevive



# ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# CMD ["/opt/GoRevive/bin/gorevive"]
CMD ["/bin/bash"]
