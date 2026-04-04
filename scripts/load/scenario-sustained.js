import http from 'k6/http';
import { check } from 'k6';

const base = __ENV.BASE_URL || 'http://127.0.0.1:8080';

// Override with e.g. SUSTAIN_VUS=60 SUSTAIN_DURATION=10m
const vus = Number(__ENV.SUSTAIN_VUS || 40);
const duration = __ENV.SUSTAIN_DURATION || '5m';

export const options = {
  scenarios: {
    sustained: {
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
  check(res, { '200': (r) => r.status === 200 });

  // Optional: hit a known id from seed data
  http.get(`${base}/albums/1`, {
    headers: { 'X-API-Key': __ENV.API_KEY },
  });
}
