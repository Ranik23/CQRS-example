import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export let options = {
  scenarios: {
    valid_orders: {
      executor: 'constant-vus',
      vus: 20,
      duration: '30s',
      exec: 'validOrders',
    },
    invalid_orders: {
      executor: 'constant-vus',
      vus: 10,
      duration: '30s',
      exec: 'invalidOrders',
    },
  },
};

const url = 'http://localhost:8080/orders';

const successCounter = new Counter('successful_orders');
const failedCounter = new Counter('failed_orders');
const invalidCounter = new Counter('invalid_orders');

const sampleItems = [
  "T-Shirt", "Jeans", "Jacket", "Sneakers", "Cap", "Backpack", "Gloves", "Scarf"
];

function generateOrder(valid = true) {
  const items = [];
  const count = randomIntBetween(1, 5);
  for (let i = 0; i < count; i++) {
    items.push({
      name: sampleItems[randomIntBetween(0, sampleItems.length - 1)],
      price: valid ? randomIntBetween(1000, 10000) : 0,
    });
  }
  return {
    user_id: valid ? randomIntBetween(1, 10000) : 0,
    items,
  };
}

export function validOrders() {
  const order = generateOrder(true);
  const res = http.post(url, JSON.stringify(order), {
    headers: { 'Content-Type': 'application/json' },
  });

  const success = check(res, {
    'valid: status is 200 or 201': (r) => r.status === 200 || r.status === 201,
  });

  if (success) {
    successCounter.add(1);
  } else {
    failedCounter.add(1);
    console.error(`❌ Failed valid order: ${JSON.stringify(order)} | response: ${res.status} - ${res.body}`);
  }

  sleep(0.2);
}

export function invalidOrders() {
  const order = generateOrder(false);
  const res = http.post(url, JSON.stringify(order), {
    headers: { 'Content-Type': 'application/json' },
  });

  const validFailure = check(res, {
    'invalid: status is 400': (r) => r.status === 400,
  });

  if (validFailure) {
    invalidCounter.add(1);
  } else {
    failedCounter.add(1);
    console.error(`❌ Unexpected success for invalid order: ${JSON.stringify(order)} | status: ${res.status}`);
  }

  sleep(0.2);
}
