# Builder Image
FROM golang:latest

WORKDIR /usr/src

COPY go.mod go.sum /usr/src/
RUN go mod download

COPY . /usr/src/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" ./cmd/static-path-provisioner

# Artifacts Image
FROM alpine:latest

COPY --from=0 /usr/src/static-path-provisioner /usr/bin/
CMD ["static-path-provisioner"]
