import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export let options = {
  vus: 50,
  duration: '30s',
};

const url = 'http://localhost:8080/orders';

const successCounter = new Counter('successful_orders');
const failedCounter = new Counter('failed_orders');

const sampleItems = [
  "T-Shirt", "Jeans", "Jacket", "Sneakers", "Cap", "Backpack", "Gloves", "Scarf"
];

function generateOrder() {
  const items = [];
  const count = randomIntBetween(1, 5);
  for (let i = 0; i < count; i++) {
    items.push({
      name: sampleItems[randomIntBetween(0, sampleItems.length - 1)],
      price: randomIntBetween(1000, 10000),
    });
  }
  return {
    user_id: randomIntBetween(1, 10000),
    items,
  };
}

export default function () {
  const order = generateOrder();
  const res = http.post(url, JSON.stringify(order), {
    headers: { 'Content-Type': 'application/json' },
  });

  const success = check(res, {
    'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
  });

  if (success) {
    successCounter.add(1);
  } else {
    failedCounter.add(1);
    console.error(`❌ Failed order: ${JSON.stringify(order)} | response: ${res.status} - ${res.body}`);
  }

  sleep(0.2);
}
