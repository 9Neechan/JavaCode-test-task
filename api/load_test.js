import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 1000 }, // За 10 сек выйти на 1000 RPS
    { duration: '30s', target: 1000 }, // Удерживать нагрузку 1000 RPS в течение 30 сек
    { duration: '10s', target: 0 },    // Постепенное снижение нагрузки
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% запросов должны быть быстрее 500мс
  },
};

export default function () {
  let walletId = "7ff05ab9-80d5-40d0-8037-7133da806e49"; // Фиксированный UUID кошелька
  let url = `http://localhost:8080/api/v1/wallets/${walletId}`;
  
  let res = http.get(url);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500
  });
}



// запуск теста
// k6 run load_test.js