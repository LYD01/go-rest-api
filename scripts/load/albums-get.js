import http from 'k6/http';
import { check } from 'k6';

const base = __ENV.BASE_URL || 'http://127.0.0.1:8080';

export const options = {
  vus: 10,
  duration: '30s',
};

export function setup() {
  if (!__ENV.API_KEY) {
    throw new Error('Set API_KEY to match the server (same as in .env)');
  }
}

export default function () {
  const res = http.get(`${base}/albums`, {
    headers: { 'X-API-Key': __ENV.API_KEY },
  });
  check(res, {
    'status 200': (r) => r.status === 200,
  });
}
