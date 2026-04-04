import http from 'k6/http';
import { check } from 'k6';

const base = __ENV.BASE_URL || 'http://127.0.0.1:8080';

export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 100 },
    { duration: '30s', target: 0 },
  ],
};

export function setup() {
  if (!__ENV.API_KEY) {
    throw new Error('Set API_KEY');
  }
}

export default function () {
  const res = http.get(`${base}/albums`, {
    headers: { 'X-API-Key': __ENV.API_KEY },
  });
  check(res, { '200': (r) => r.status === 200 });
}
