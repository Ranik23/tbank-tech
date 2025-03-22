import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 20 },
    { duration: '20s', target: 100 },
    { duration: '10s', target: 50 },
    { duration: '10s', target: 0},
  ],
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:9091';

export default function () {
  let userId = Math.floor(Math.random() * 100000);
  let token = 'test-token';

  let registerRes = http.post(`${BASE_URL}/user/${userId}`, JSON.stringify({
    name: `user${userId}`,
    token: token,
  }), { headers: { 'Content-Type': 'application/json' } });

  check(registerRes, {
    'RegisterUser HTTP 200 or 409': (r) => r.status === 200 || r.status === 409,
  });

  sleep(1);


  let getLinksRes = http.get(`${BASE_URL}/users/${userId}/links`);

  check(getLinksRes, {
    'GetLinks HTTP 200': (r) => r.status === 200,  
  })


  sleep(1);

  let deleteRes = http.del(`${BASE_URL}/user/${userId}`);

  check(deleteRes, {
    'DeleteUser HTTP 200 or 404': (r) => r.status === 200 || r.status === 404,
  });

  sleep(1);
}
