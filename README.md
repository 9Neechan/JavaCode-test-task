Описание решения

wait-for.sh: https://github.com/eficode/wait-for/releases 

db scheme: https://dbdiagram.io/d/67ab470b263d6cf9a0c45391

Запуск сервиса:

```docker-compose up```

Нагрузочное тестирование ks6 1000rps

```k6 run api/load_test2.js```



minikube start --driver=docker

minikube stop
minikube delete
