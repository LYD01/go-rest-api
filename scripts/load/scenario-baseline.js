import http from 'k6/http';
import { check } from 'k6';

const base = __ENV.BASE_URL || 'http://127.0.0.1:8080';

export const options = {
  vus: 1,
  duration: '30s',
};

export function setup() {
  if (!__ENV.API_KEY) {
    throw new Error('Set API_KEY for album checks');
  }
}

export default function () {
  const h = http.get(`${base}/healthz`);
  check(h, { 'health 200': (r) => r.status === 200 });

  const a = http.get(`${base}/albums`, {
    headers: { 'X-API-Key': __ENV.API_KEY },
  });
  check(a, { 'albums 200': (r) => r.status === 200 });
}
