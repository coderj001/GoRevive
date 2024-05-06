# Build Stage
FROM coderj001/gorevive:1.13 AS build-stage

LABEL app="build-GoRevive"
LABEL REPO="https://github.com/coderj001/GoRevive"

ENV PROJPATH=/go/src/github.com/coderj001/GoRevive

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/coderj001/GoRevive
WORKDIR /go/src/github.com/coderj001/GoRevive

RUN make build-alpine

# Final Stage
FROM coderj001/gorevive

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/coderj001/GoRevive"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/GoRevive/bin

WORKDIR /opt/GoRevive/bin

COPY --from=build-stage /go/src/github.com/coderj001/GoRevive/bin/GoRevive /opt/GoRevive/bin/
RUN chmod +x /opt/GoRevive/bin/GoRevive

# Create appuser
RUN adduser -D -g '' GoRevive
USER GoRevive

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/GoRevive/bin/GoRevive"]
