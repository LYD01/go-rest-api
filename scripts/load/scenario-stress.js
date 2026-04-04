import http from 'k6/http';
import { check } from 'k6';

const base = __ENV.BASE_URL || 'http://127.0.0.1:8080';

// Intentionally heavy — run only against localhost you control.
const vus = Number(__ENV.STRESS_VUS || 200);
const duration = __ENV.STRESS_DURATION || '2m';

export const options = {
  scenarios: {
    stress: {
      executor: 'constant-vus',
      vus,
      duration,
    },
  },
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
  check(res, {
    '200': (r) => r.status === 200,
    '401': (r) => r.status === 401,
  });
}
