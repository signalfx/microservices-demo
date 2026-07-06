FROM eclipse-temurin:21-jdk-alpine-3.23 AS builder

RUN apk add --no-cache curl gcompat

WORKDIR /app

COPY build.gradle settings.gradle gradle.lockfile gradlew ./
COPY gradle gradle
RUN sed -i 's/\r$//' gradlew && chmod +x gradlew && \
    ./gradlew --no-daemon dependencies

COPY src src
RUN ./gradlew --no-daemon check installDist

ARG SPLUNK_OTEL_JAVA_VERSION=v2.29.0
ARG SPLUNK_OTEL_JAVA_SHA256=54e4d1943ef9903889053184ffc09969e2a7b8c308b29812c855caa54f02aaf6
RUN curl -fsSL "https://github.com/signalfx/splunk-otel-java/releases/download/${SPLUNK_OTEL_JAVA_VERSION}/splunk-otel-javaagent.jar" \
        -o /tmp/splunk-otel-javaagent.jar && \
    echo "${SPLUNK_OTEL_JAVA_SHA256}  /tmp/splunk-otel-javaagent.jar" | sha256sum -c -

ARG GRPC_HEALTH_PROBE_VERSION=v0.4.52
ARG GRPC_HEALTH_PROBE_SHA256=48a2103bb33e3d051dff9b909e340f3f2ec9e07abc467e4ef5008e88f742ae4e
RUN curl -fsSL "https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64" \
        -o /tmp/grpc_health_probe && \
    echo "${GRPC_HEALTH_PROBE_SHA256}  /tmp/grpc_health_probe" | sha256sum -c - && \
    chmod +x /tmp/grpc_health_probe

FROM eclipse-temurin:21-jre-noble

RUN groupadd --system app && \
    useradd --system --gid app --home-dir /app app && \
    mkdir -p /opt/sfx

COPY --from=builder /tmp/splunk-otel-javaagent.jar /opt/sfx/splunk-otel-javaagent-all.jar
COPY --from=builder /tmp/grpc_health_probe /bin/grpc_health_probe

WORKDIR /app
COPY --from=builder --chown=app:app /app/build/install/AdService/ ./

USER app
EXPOSE 9555
ENTRYPOINT ["/app/bin/AdService"]
