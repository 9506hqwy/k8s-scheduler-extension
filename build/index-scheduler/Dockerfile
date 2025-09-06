FROM golang:1.24-bookworm as build

WORKDIR /src

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN --mount=type=bind,target=/src,Z \
    go build -o /index-scheduler ./cmd/index-scheduler

FROM gcr.io/distroless/static-debian12:latest

COPY --from=build /index-scheduler /

ENTRYPOINT ["/index-scheduler"]
