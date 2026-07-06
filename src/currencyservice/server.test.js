'use strict';

const assert = require('node:assert/strict');
const test = require('node:test');

const { _carry, convert } = require('./server');

test('carry normalizes fractional units and nanos', () => {
  assert.deepEqual(_carry({ units: 1.5, nanos: 750000000 }), {
    units: 2,
    nanos: 250000000
  });
});

test('convert returns a normalized target currency amount', async () => {
  const result = await new Promise((resolve, reject) => {
    convert({
      request: {
        from: { currency_code: 'EUR', units: 10, nanos: 0 },
        to_code: 'USD'
      }
    }, (err, response) => {
      if (err) {
        reject(err);
      } else {
        resolve(response);
      }
    });
  });

  assert.equal(result.currency_code, 'USD');
  assert.equal(Number.isInteger(result.units), true);
  assert.equal(Number.isInteger(result.nanos), true);
});
