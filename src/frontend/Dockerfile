FROM golang:1.26.4-alpine3.24 AS builder
RUN apk add --no-cache ca-certificates git

ARG GRPC_HEALTH_PROBE_VERSION=v0.4.52
ARG GRPC_HEALTH_PROBE_SHA256=48a2103bb33e3d051dff9b909e340f3f2ec9e07abc467e4ef5008e88f742ae4e
RUN wget -qO /bin/grpc_health_probe "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64" && \
    echo "${GRPC_HEALTH_PROBE_SHA256}  /bin/grpc_health_probe" | sha256sum -c - && \
    chmod +x /bin/grpc_health_probe

WORKDIR /app
COPY . .
RUN go build -o /go/bin/frontend .

FROM alpine:3.24 AS release
RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S -G app app
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
WORKDIR /frontend
COPY --from=builder /go/bin/frontend /frontend/server
COPY ./templates ./templates
COPY ./static ./static
USER app
EXPOSE 8080
ENTRYPOINT ["/frontend/server"]
