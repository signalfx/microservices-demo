import os
import unittest

os.environ.setdefault('SIGNALFX_ENDPOINT_URL', 'http://127.0.0.1:9411/api/v2/spans')

import demo_pb2
from recommendation_server import RecommendationService


class ProductCatalogStub:
    def ListProducts(self, request):
        return demo_pb2.ListProductsResponse(products=[
            demo_pb2.Product(id='1'),
            demo_pb2.Product(id='2'),
            demo_pb2.Product(id='3'),
        ])


class RecommendationServiceTest(unittest.TestCase):
    def test_recommendations_exclude_products_already_in_cart(self):
        service = RecommendationService(ProductCatalogStub())
        request = demo_pb2.ListRecommendationsRequest(product_ids=['2'])

        response = service.ListRecommendations(request, None)

        self.assertCountEqual(response.product_ids, ['1', '3'])


if __name__ == '__main__':
    unittest.main()
