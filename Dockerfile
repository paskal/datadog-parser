# Build
FROM umputun/baseimage:buildgo-latest as build

ARG SKIP_TEST

ADD . /app
WORKDIR /app

# run tests
RUN \
    if [ -z "$SKIP_TEST" ] ; then \
        go test -timeout=30s  ./... && \
        golangci-lint run --config .golangci.yml ./... ; \
    else echo "skip tests and linter" ; fi

RUN go build -o datadog-parser ./app

# Run
FROM alpine:latest

COPY --from=build /app/datadog-parser /app/

WORKDIR /app

ENTRYPOINT ["/app/datadog-parser"]
