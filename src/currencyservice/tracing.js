const { NodeTracerProvider } = require('@opentelemetry/sdk-trace-node');
const { B3InjectEncoding, B3Propagator } = require('@opentelemetry/propagator-b3');
const { ZipkinExporter } = require('@opentelemetry/exporter-zipkin');
const { GrpcInstrumentation } = require('@opentelemetry/instrumentation-grpc');
const { BatchSpanProcessor } = require('@opentelemetry/sdk-trace-base');
const { registerInstrumentations } = require('@opentelemetry/instrumentation');

const exporter = new ZipkinExporter({
  serviceName: 'currencyservice',
  url: process.env.SIGNALFX_ENDPOINT_URL
});

const provider = new NodeTracerProvider({
  spanProcessors: [new BatchSpanProcessor(exporter)]
});
provider.register({
  propagator: new B3Propagator({
    injectEncoding: B3InjectEncoding.MULTI_HEADER
  })
});

registerInstrumentations({
  instrumentations: [
    new GrpcInstrumentation()
  ]
});
