import http from 'k6/http';
import { check } from 'k6';
import { sleep } from 'k6';

export let options = {
  vus: 10, // Number of virtual users
  duration: '30s', // Duration of the test
};

export default function () {
  let res = http.get('http://localhost:8888/ping');
  check(res, {
    'is status 200': (r) => r.status === 200,
  });

  // Optional: Add a small delay between requests to simulate more realistic user behavior
  sleep(1);
}