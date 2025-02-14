## Описание решения

Представлены 2 решения: в папках api, api_nobalance

api: нагрузка балансируется с помощью кеширования redis и очередей rabbitmq

api_nobalance: нагрузка не балансируется, для нее написаны тесты

если в main.go поменять строку инициализации сервера, то можно его запустить и проверить работу тестов

### Запуск сервиса:

Докер:

```docker-compose up```

Локально:

```make build```

### Тестовые запросы

GET 

```http://localhost:8080/api/v1/wallets/7ff05ab9-80d5-40d0-8037-7133da806e49```

POST

```http://localhost:8080/api/v1/wallet```

html
```
{
  "amount": 100,
  "wallet_uuid": "7ff05ab9-80d5-40d0-8037-7133da806e49",
  "operation_type": "DEPOSIT"
}
```

### Нагрузочное тестирование ks6 - не проходит (мб моему ноутбуку не хватает мощности)

```k6 run api/load_test3.js```

часть GET запросов завершаются с ошибкой 404 - не знаю как исправить, нет опыта в обеспечении 1000rps

Можно было еще балансировать нагрузку с помощью kubernetes/nginx

### Краткое описание логики запросов

GET: (api/wallet_redis_cache.go)

- Сначала попытаться получить баланс из кэша.
- Если в кэше нет данных, запросить их из БД.
- После запроса записать данные в Redis и вернуть пользователю.

POST: (api/wallet_rabbitmq.go)

- Ставит обновления в очередь RabbitMQ, реализовывала в первый раз
- Записывает обновления в кэш - кэширует на 5 секунд
- Обновление реализовано через транзакцию (получение баланса + обновление или ошибка если недоставточно средств)

#### Доп материалы

wait-for.sh: https://github.com/eficode/wait-for/releases 

db scheme: https://dbdiagram.io/d/67ab470b263d6cf9a0c45391