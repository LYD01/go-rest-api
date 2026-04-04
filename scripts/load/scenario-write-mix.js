import http from 'k6/http';
import { check } from 'k6';

const base = __ENV.BASE_URL || 'http://127.0.0.1:8080';

export const options = {
  vus: 30,
  duration: '2m',
};

const headers = (key) => ({
  headers: {
    'X-API-Key': key,
    'Content-Type': 'application/json',
  },
});

export function setup() {
  if (!__ENV.API_KEY) {
    throw new Error('Set API_KEY');
  }
}

export default function () {
  const key = __ENV.API_KEY;
  const getRes = http.get(`${base}/albums`, headers(key));
  check(getRes, { 'GET 200': (r) => r.status === 200 });

  const body = JSON.stringify({
    title: `load-${__VU}-${__ITER}`,
    artist: 'k6',
    price: 9.99,
  });
  const postRes = http.post(`${base}/albums`, body, headers(key));
  check(postRes, {
    'POST 201': (r) => r.status === 201,
  });
}
