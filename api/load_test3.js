import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  scenarios: {
    get_wallet: {
      executor: 'constant-arrival-rate',
      rate: 500, // запросов в секунду
      timeUnit: '1s',
      duration: '50s', 
      preAllocatedVUs: 500,
      maxVUs: 1000,
    },
    update_wallet: {
      executor: 'constant-arrival-rate',
      rate: 500,
      timeUnit: '1s',
      duration: '50s',
      preAllocatedVUs: 500,
      maxVUs: 1000,
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
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
  
  // Отправка в очередь RabbitMQ через API
  //let res = http.post(`${BASE_URL}/wallet/queue`, payload, { headers });
  let res = http.post(`${BASE_URL}/wallet`, payload, { headers });

  check(res, {
    'POST status is 202': (r) => r.status === 202, // Теперь ждем 202, так как обработка асинхронная
  });

  if (res.status === 202) {
    // Пытаемся получить статус обновления через polling
    let maxRetries = 5;
    let attempt = 0;
    
    while (attempt < maxRetries) {
      sleep(0.5); // Подождем 500 мс перед следующим запросом
      let checkRes = http.get(`${BASE_URL}/wallets/${WALLET_ID}/status`);

      if (checkRes.status === 200) {
        let jsonResponse = checkRes.json();
        if (jsonResponse.status === "COMPLETED") {
          check(checkRes, { 'Update completed successfully': () => true });
          return;
        }
      }
      attempt++;
    }
  }
}

// Разделяем выполнение тестов по сценариям
export default function () {
  if (__VU % 2 === 0) {
    getWallet();
  } else {
    updateWallet();
  }
}
