#!/usr/bin/python
#
# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import random
import time
from concurrent import futures

import grpc

from opentelemetry import trace
from opentelemetry.propagate import set_global_textmap
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.exporter.zipkin.json import ZipkinExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.propagators.b3 import B3MultiFormat
from opentelemetry.instrumentation.grpc import GrpcInstrumentorServer
from opentelemetry.instrumentation.grpc import GrpcInstrumentorClient

import demo_pb2
import demo_pb2_grpc
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc
from grpc_health.v1 import health

from logger import getJSONLogger
logger = getJSONLogger('recommendationservice-server')

zipkin_exporter = ZipkinExporter(endpoint=os.environ['SIGNALFX_ENDPOINT_URL'])
span_processor = BatchSpanProcessor(zipkin_exporter)

set_global_textmap(B3MultiFormat())
provider = TracerProvider(
    resource=Resource.create({'service.name': 'recommendationservice'})
)
provider.add_span_processor(span_processor)
trace.set_tracer_provider(provider)
tracer = trace.get_tracer(__name__)

instrumentor = GrpcInstrumentorClient()
instrumentor.instrument()
grpc_server_instrumentor = GrpcInstrumentorServer()
grpc_server_instrumentor.instrument()


class RecommendationService(demo_pb2_grpc.RecommendationServiceServicer):
    def __init__(self, product_catalog_stub):
        self.product_catalog_stub = product_catalog_stub

    def ListRecommendations(self, request, context):
        max_responses = 5
        # fetch list of products from product catalog stub
        cat_response = self.product_catalog_stub.ListProducts(demo_pb2.Empty())
        product_ids = [x.id for x in cat_response.products]
        filtered_products = list(set(product_ids)-set(request.product_ids))
        num_products = len(filtered_products)
        num_return = min(max_responses, num_products)
        # sample list of indicies to return
        indices = random.sample(range(num_products), num_return)
        # fetch product ids from indices
        prod_list = [filtered_products[i] for i in indices]
        logger.info("[Recv ListRecommendations] product_ids={}".format(prod_list))
        # build and return response
        response = demo_pb2.ListRecommendationsResponse()
        response.product_ids.extend(prod_list)
        return response

if __name__ == "__main__":
    logger.info("initializing recommendationservice")

    logger.info("Tracing enabled.")

    port = os.environ.get('PORT', "8080")
    catalog_addr = os.environ.get('PRODUCT_CATALOG_SERVICE_ADDR', '')
    if catalog_addr == "":
        raise Exception('PRODUCT_CATALOG_SERVICE_ADDR environment variable not set')
    logger.info("product catalog address: " + catalog_addr)
    ca_path = os.environ.get('GRPC_TLS_CA', '')
    if ca_path == "":
        raise Exception('GRPC_TLS_CA environment variable not set')
    with open(ca_path, 'rb') as ca_file:
        channel_credentials = grpc.ssl_channel_credentials(
            root_certificates=ca_file.read()
        )
    channel = grpc.secure_channel(catalog_addr, channel_credentials)
    product_catalog_stub = demo_pb2_grpc.ProductCatalogServiceStub(channel)

    # create gRPC server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # add class to gRPC server
    service = RecommendationService(product_catalog_stub)
    health_service = health.HealthServicer()
    health_service.set('', health_pb2.HealthCheckResponse.SERVING)
    demo_pb2_grpc.add_RecommendationServiceServicer_to_server(service, server)
    health_pb2_grpc.add_HealthServicer_to_server(health_service, server)

    # start server
    logger.info("listening on port: " + port)
    server.add_insecure_port('[::]:'+port)
    server.start()

    # keep alive
    try:
         while True:
            time.sleep(10000)
    except KeyboardInterrupt:
            server.stop(0)
