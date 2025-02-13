import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  scenarios: {
    get_wallet: {
      executor: 'constant-arrival-rate',
      rate: 500, // запросов в секунду
      timeUnit: '1s',
      duration: '50s', // Общая длительность теста
      preAllocatedVUs: 500, // Начальное число виртуальных пользователей
      maxVUs: 750, // Максимальное число виртуальных пользователей
    },
    update_wallet: {
      executor: 'constant-arrival-rate',
      rate: 500, // запросов в секунду
      timeUnit: '1s',
      duration: '50s',
      preAllocatedVUs: 500,
      maxVUs: 750,
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% запросов должны быть <500ms
    http_req_failed: ['rate<0.01'],   // Ошибок должно быть <1%
  },
};

const WALLET_ID = "7ff05ab9-80d5-40d0-8037-7133da806e49";
const BASE_URL = "http://localhost:8080/api/v1";

export function getWallet() {
  let res = http.get(`${BASE_URL}/wallets/${WALLET_ID}`);

  check(res, {
    'GET status is 200': (r) => r.status === 200,
    'GET response time < 500ms': (r) => r.timings.duration < 500
  });
}

export function updateWallet() {
  let payload = JSON.stringify({
    amount: 1,
    wallet_uuid: WALLET_ID,
    operation_type: "DEPOSIT"
  });

  let headers = { 'Content-Type': 'application/json' };
  let res = http.post(`${BASE_URL}/wallet`, payload, { headers });

  check(res, {
    'POST status is 200': (r) => r.status === 200,
    'POST response time < 500ms': (r) => r.timings.duration < 500
  });
}

// Разделяем выполнение тестов по сценариям
export default function () {
  if (__VU % 2 === 0) {
    getWallet();
  } else {
    updateWallet();
  }
}
