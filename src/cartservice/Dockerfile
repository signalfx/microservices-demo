ARG NET_VERSION=10.0

FROM mcr.microsoft.com/dotnet/sdk:${NET_VERSION} AS builder

WORKDIR /source
COPY cartservice.csproj packages.lock.json ./
RUN dotnet restore cartservice.csproj --locked-mode

COPY . .
RUN dotnet build cartservice.csproj -c Release --no-restore && \
    dotnet publish cartservice.csproj -c Release --no-self-contained --no-restore -o /cartservice

FROM mcr.microsoft.com/dotnet/aspnet:${NET_VERSION} AS instrumentation

ADD --checksum=sha256:0abfc82eb0408838636e9b1e0557aefc0a4c81ccc441b6b39807120f9800015b \
    https://github.com/signalfx/splunk-otel-dotnet/releases/download/v1.14.0/splunk-opentelemetry-dotnet-linux-glibc-x64.zip \
    /tmp/splunk-otel-dotnet.zip
RUN apt-get update && \
    apt-get install -y --no-install-recommends unzip && \
    mkdir -p /.splunk-otel-dotnet && \
    unzip -q /tmp/splunk-otel-dotnet.zip -d /.splunk-otel-dotnet

FROM mcr.microsoft.com/dotnet/aspnet:${NET_VERSION}

ADD --checksum=sha256:48a2103bb33e3d051dff9b909e340f3f2ec9e07abc467e4ef5008e88f742ae4e \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.52/grpc_health_probe-linux-amd64 \
    /bin/grpc_health_probe
RUN chmod +x /bin/grpc_health_probe && \
    groupadd --system cart && \
    useradd --system --gid cart --home-dir /app cart

COPY --from=instrumentation /.splunk-otel-dotnet /.splunk-otel-dotnet

WORKDIR /app
COPY --from=builder --chown=cart:cart /cartservice ./

ENV CORECLR_ENABLE_PROFILING=1 \
    CORECLR_PROFILER={918728DD-259F-4A6A-AC2B-B85E1B658318} \
    CORECLR_PROFILER_PATH=/.splunk-otel-dotnet/linux-x64/OpenTelemetry.AutoInstrumentation.Native.so \
    DOTNET_ADDITIONAL_DEPS=/.splunk-otel-dotnet/AdditionalDeps \
    DOTNET_SHARED_STORE=/.splunk-otel-dotnet/store \
    DOTNET_STARTUP_HOOKS=/.splunk-otel-dotnet/net/OpenTelemetry.AutoInstrumentation.StartupHook.dll \
    OTEL_DOTNET_AUTO_HOME=/.splunk-otel-dotnet \
    OTEL_DOTNET_AUTO_PLUGINS="Splunk.OpenTelemetry.AutoInstrumentation.Plugin, Splunk.OpenTelemetry.AutoInstrumentation"

USER cart
EXPOSE 7070
ENTRYPOINT ["./cartservice", "start"]
