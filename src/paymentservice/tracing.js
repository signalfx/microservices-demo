const { NodeTracerProvider } = require('@opentelemetry/sdk-trace-node');
const { B3InjectEncoding, B3Propagator } = require('@opentelemetry/propagator-b3');
const { ZipkinExporter } = require('@opentelemetry/exporter-zipkin');
const { GrpcInstrumentation } = require('@opentelemetry/instrumentation-grpc');
const { BatchSpanProcessor, ConsoleSpanExporter } = require('@opentelemetry/sdk-trace-base');
const { registerInstrumentations } = require('@opentelemetry/instrumentation');

const spanProcessors = [new BatchSpanProcessor(new ZipkinExporter({
  serviceName: 'paymentservice',
  url: process.env.SIGNALFX_ENDPOINT_URL
}))];

if (process.env.CONSOLE_SPAN === 'true') {
  spanProcessors.push(new BatchSpanProcessor(new ConsoleSpanExporter()));
}

const provider = new NodeTracerProvider({ spanProcessors });
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
