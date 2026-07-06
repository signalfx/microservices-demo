FROM node:24.18.0-alpine3.24 AS builder

ARG GRPC_HEALTH_PROBE_VERSION=v0.4.52
ARG GRPC_HEALTH_PROBE_SHA256=48a2103bb33e3d051dff9b909e340f3f2ec9e07abc467e4ef5008e88f742ae4e
RUN wget -qO /bin/grpc_health_probe "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64" && \
    echo "${GRPC_HEALTH_PROBE_SHA256}  /bin/grpc_health_probe" | sha256sum -c - && \
    chmod +x /bin/grpc_health_probe

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm ci --omit=dev --ignore-scripts && npm cache clean --force

FROM node:24.18.0-alpine3.24

RUN rm -rf /usr/local/lib/node_modules

WORKDIR /usr/src/app

COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=builder --chown=node:node /usr/src/app/node_modules ./node_modules

COPY --chown=node:node . .

ENV NODE_ENV=production
USER node
EXPOSE 7000

ENTRYPOINT [ "node", "-r", "./tracing.js", "server.js" ]
