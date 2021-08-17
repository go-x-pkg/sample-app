FROM debian:11 as base

ENV GO_VERSION="1.13.4"

RUN \
  apt-get update && \
  apt-get install -y \
  curl \
  git-core

# golang
WORKDIR /opt
RUN \
  curl -OL "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" && \
  mv ./go*tar.gz ./go.tar.gz && \
  tar -xf ./go.tar.gz && \
  rm -rf ./go.tar.gz

ENV GOROOT="/opt/go"
ENV PATH="${GOROOT}/bin:${PATH}"

WORKDIR /usr/local/sample-app

COPY ./go.mod ./
RUN go mod download

COPY ./cmd/ ./cmd
COPY ./internal/ ./internal/
COPY ./pkg/ ./pkg/
COPY ./appversion/ ./appversion/

COPY ./hack/build-sample-app-server.sh ./hack/build-sample-app-server.sh

# -----------------------------------------------------------------------------
FROM base as server-build

RUN /bin/bash ./hack/build-sample-app-server.sh
RUN cp /usr/local/sample-app/bin/sample-app-server /usr/local/bin/sample-app-server

# -----------------------------------------------------------------------------
FROM server-build as server-run

COPY --from=server-build /usr/local/bin/sample-app-server /usr/local/bin/sample-app-server

CMD ["--help"]
ENTRYPOINT ["app-server"]
