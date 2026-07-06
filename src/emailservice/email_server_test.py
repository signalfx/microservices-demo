import os
import unittest

os.environ.setdefault('SIGNALFX_ENDPOINT_URL', 'http://127.0.0.1:9411/api/v2/spans')

import demo_pb2
from email_server import DummyEmailService


class DummyEmailServiceTest(unittest.TestCase):
  def test_send_order_confirmation_returns_empty_response(self):
    request = demo_pb2.SendOrderConfirmationRequest(email='test@example.com')

    response = DummyEmailService().SendOrderConfirmation(request, None)

    self.assertIsInstance(response, demo_pb2.Empty)


if __name__ == '__main__':
  unittest.main()
