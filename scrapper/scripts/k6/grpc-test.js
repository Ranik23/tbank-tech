import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '20s', target: 95 },
    { duration: '10s', target: 0 },
  ],
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8082'; 

const client = new grpc.Client();
  client.load(['/home/anton/tbank-tech/scrapper/api/third_party/google/api'], 'annotations.proto', 'http.proto');

  client.load(['/home/anton/tbank-tech/scrapper/api/proto'], 'scrapper.proto')


export default function () {
  let userId = Math.floor(Math.random() * 100000);
  let token = 'test-token';

  // 1. Регистрация пользователя через gRPC
  let registerRes = client.invoke('scrapper.Scrapper.RegisterUser', {
    tg_user_id: userId,
    name: `user${userId}`,
    token: token,
  });

  check(registerRes, {
    'RegisterUser gRPC success': (r) => r.status === 0 || r.status === 6,
  });

  sleep(1);

  // 2. Получение ссылок через gRPC
  let getLinksRes = client.invoke('scrapper.Scrapper.GetLinks', { tg_user_id: userId });

  check(getLinksRes, {
    'GetLinks gRPC success': (r) => r.status === 0,
  });

  sleep(1);

  let deleteRes = client.invoke('scrapper.Scrapper.DeleteUser', { tg_user_id: userId });

  check(deleteRes, {
    'DeleteUser gRPC success': (r) => r.status === 0 || r.status === 5,
  });

  sleep(1);
}

