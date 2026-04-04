import http from 'k6/http';
import { check } from 'k6';

const base = __ENV.BASE_URL || 'http://127.0.0.1:8080';

export const options = {
  vus: 5,
  duration: '30s',
};

export default function () {
  const res = http.get(`${base}/healthz`);
  check(res, {
    'status 200': (r) => r.status === 200,
    'json status ok': (r) => {
      try {
        const b = r.json();
        return b && b.status === 'ok';
      } catch {
        return false;
      }
    },
  });
}
