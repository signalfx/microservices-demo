// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/contrib/propagators/b3"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	pb "github.com/signalfx/microservices-demo/src/shippingservice/genproto"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	defaultPort = "50051"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	logger.Out = os.Stdout
}

func main() {
	if os.Getenv("DISABLE_TRACING") == "" {
		logger.Info("Tracing enabled.")
		stopTracing := initTracing()
		defer stopTracing()
	} else {
		logger.Info("Tracing disabled.")
	}

	if os.Getenv("DISABLE_PROFILER") == "" {
		logger.Info("Profiling enabled.")
		go initProfiling("shippingservice", "1.0.0")
	} else {
		logger.Info("Profiling disabled.")
	}

	port := defaultPort
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}
	port = fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	var srv *grpc.Server
	if os.Getenv("DISABLE_STATS") == "" {
		logger.Info("Stats enabled.")
		srv = grpc.NewServer(
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)
	} else {
		logger.Info("Stats disabled.")
		srv = grpc.NewServer()
	}
	svc := &server{}
	pb.RegisterShippingServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)
	logger.Infof("Shipping Service listening on port %s", port)

	// Register reflection service on gRPC server.
	reflection.Register(srv)
	if err := srv.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

// server controls RPC service responses.
type server struct{}

// Check is for health checking.
func (s *server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *server) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

// GetQuote produces a shipping quote (cost) in USD.
func (s *server) GetQuote(ctx context.Context, in *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	log := logger.WithFields(getTraceLogFields(ctx))
	log.Info("[GetQuote] received request")
	defer log.Info("[GetQuote] completed request")

	// 1. Our quote system requires the total number of items to be shipped.
	count := 0
	for _, item := range in.Items {
		count += int(item.Quantity)
	}

	// 2. Generate a quote based on the total number of items to be shipped.
	quote := CreateQuoteFromCount(count)

	// 3. Generate a response.
	return &pb.GetQuoteResponse{
		CostUsd: &pb.Money{
			CurrencyCode: "USD",
			Units:        int64(quote.Dollars),
			Nanos:        int32(quote.Cents * 10000000)},
	}, nil

}

// ShipOrder mocks that the requested items will be shipped.
// It supplies a tracking ID for notional lookup of shipment delivery status.
func (s *server) ShipOrder(ctx context.Context, in *pb.ShipOrderRequest) (*pb.ShipOrderResponse, error) {
	log := logger.WithFields(getTraceLogFields(ctx))
	log.Info("[ShipOrder] received request")
	defer log.Info("[ShipOrder] completed request")
	// 1. Create a Tracking ID
	baseAddress := fmt.Sprintf("%s, %s, %s", in.Address.StreetAddress, in.Address.City, in.Address.State)
	id := CreateTrackingId(baseAddress)

	// 2. Generate a response.
	return &pb.ShipOrderResponse{
		TrackingId: id,
	}, nil
}

/*
func initTracing() func() {
	sdk, err := distro.Run(distro.WithServiceName("shippingservice"))
	if err != nil {
		panic(err)
	}

	logger.Info("otel initialization completed.")
	return func() {
		if err := sdk.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}
}
*/

func initTracing() func() {
	endpoint := os.Getenv("OTEL_EXPORTER_JAEGER_ENDPOINT")
	if endpoint == "" { endpoint = "localhost:14268/api/traces" }

	exp, err := jaeger.NewRawExporter(
				jaeger.WithCollectorEndpoint(endpoint),
				)

	if err != nil {
		fmt.Errorf("%s: %v", "failed to create exporter", err)
		os.Exit(1)
	}
	ctx := context.Background()
	res, _ := resource.New(ctx)

	traceProvider := sdktrace.NewTracerProvider(
					sdktrace.WithSampler(sdktrace.AlwaysSample()),
					sdktrace.WithResource(res),
					sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exp)),
					)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(b3.B3{})

	logger.Info("otel initialization completed.")
	return func() {
		err := traceProvider.Shutdown(ctx)
		if err != nil {
			fmt.Errorf("%s: %v", "failed to shutdown provider", err)
			os.Exit(1)
		}
		err = exp.Shutdown(ctx)
		if err != nil {
			fmt.Errorf("%s: %v", "failed to stop exporter", err)
			os.Exit(1)
		}
	}
}

func initProfiling(service, version string) {
	// TODO(ahmetb) this method is duplicated in other microservices using Go
	// since they are not sharing packages.
	for i := 1; i <= 3; i++ {
		if err := profiler.Start(profiler.Config{
			Service:        service,
			ServiceVersion: version,
			// ProjectID must be set if not running on GCP.
			// ProjectID: "my-project",
		}); err != nil {
			logger.Warnf("failed to start profiler: %+v", err)
		} else {
			logger.Info("started Stackdriver profiler")
			return
		}
		d := time.Second * 10 * time.Duration(i)
		logger.Infof("sleeping %v to retry initializing Stackdriver profiler", d)
		time.Sleep(d)
	}
	logger.Warn("could not initialize Stackdriver profiler after retrying, giving up")
}

func getTraceLogFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}
	if span := trace.SpanFromContext(ctx); span != nil {
		spanCtx := span.SpanContext()
		fields["trace_id"] = spanCtx.TraceID().String()
		fields["span_id"] = spanCtx.SpanID().String()
		fields["service.name"] = "shippingservice"
	}
	return fields
}
