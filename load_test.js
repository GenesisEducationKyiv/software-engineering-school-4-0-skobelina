import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 100 },
    { duration: '1m', target: 100 },
    { duration: '10s', target: 0 }, 
    ],
};

const API_BASE_URL = 'http://localhost:8080/api';

export default function () {
// tests rate endpoint
    let rateRes = http.get(`${API_BASE_URL}/rate`);
  check(rateRes, {
    'rate status is 200': (r) => r.status === 200,
  });

// tests subscribe endpoint
  let subscribePayload = JSON.stringify({
    email: `skobelina${__VU}@gmail.com`,
  });
  let subscribeRes = http.post(`${API_BASE_URL}/subscribe`, subscribePayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(subscribeRes, {
    'subscribe status is 200': (r) => r.status === 200 || r.status === 409,
  });

// tests notification endpoint
  let key = 'E6B3C4F7';
  let cronJobRes = http.get(`${API_BASE_URL}/cron-jobs/notifications/exchange-rates/${key}`);
  check(cronJobRes, {
    'cron job status is 200': (r) => r.status === 200 || r.status === 204, 
  });

  sleep(1);
}

