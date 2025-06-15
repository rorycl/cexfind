# build the webapp including dependencies from / and /cmd

# For information about the use of the gcr.io/distroless/static final
# container image, see the gcr.io docs at
# https://github.com/GoogleContainerTools/distroless/blob/main/base/README.md

FROM golang:1.24 AS deps

# setup module environment
WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
ADD cmd/web/*mod cmd/web/*sum ./cmd/web/
RUN cd /build/cmd/web && go mod download

# build
FROM deps as dev
ADD *go ./
ADD cmd/*go ./cmd/
ADD cmd/web/ ./cmd/web/
ADD location/ ./location/
RUN cd cmd/web && \
    CGO_ENABLED=0 GOOS=linux \
    go build -ldflags "-w -X main.docker=true" -o /build/webserver .

# install into minimal image
FROM gcr.io/distroless/base AS base
WORKDIR /
EXPOSE 8000
COPY --from=dev /build/webserver /
CMD ["/webserver", "--address", "0.0.0.0", "--port", "8000"]
