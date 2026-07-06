FROM python:3.13.14-alpine3.24 AS builder

RUN apk add --no-cache \
        build-base \
        g++ \
        musl-dev \
        libffi-dev

COPY requirements.txt .
RUN pip wheel --no-cache-dir --wheel-dir /wheels -r requirements.txt

ARG GRPC_HEALTH_PROBE_VERSION=v0.4.52
ARG GRPC_HEALTH_PROBE_SHA256=48a2103bb33e3d051dff9b909e340f3f2ec9e07abc467e4ef5008e88f742ae4e
RUN wget -qO /bin/grpc_health_probe "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64" && \
    echo "${GRPC_HEALTH_PROBE_SHA256}  /bin/grpc_health_probe" | sha256sum -c - && \
    chmod +x /bin/grpc_health_probe

FROM python:3.13.14-alpine3.24

COPY --from=builder /wheels /wheels
COPY requirements.txt .
RUN pip install --no-cache-dir --no-index --find-links=/wheels -r requirements.txt && \
    rm -rf /wheels && \
    addgroup -S app && adduser -S -G app app

COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

WORKDIR /recommendationservice

COPY --chown=app:app . .

ENV PORT=8080 \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1
USER app
EXPOSE 8080

ENTRYPOINT ["python", "recommendation_server.py"]
