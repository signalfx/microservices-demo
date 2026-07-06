'use strict';

process.env.SUCCESS_PAYMENT_SERVICE_DURATION_MILLIS = '0';

const assert = require('node:assert/strict');
const test = require('node:test');

const charge = require('./charge');

test('charge returns a UUID for a valid Visa payment', async () => {
  const response = await charge({
    amount: {
      currency_code: 'USD',
      units: 10,
      nanos: 0
    },
    credit_card: {
      credit_card_number: '4111111111111111',
      credit_card_cvv: 123,
      credit_card_expiration_year: new Date().getFullYear() + 1,
      credit_card_expiration_month: 12
    }
  });

  assert.match(
    response.transaction_id,
    /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/
  );
});
