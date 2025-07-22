# Effective Mobile test assesment

Порт приложения `8080` (указан в конфиге .env)

```sh
docker compose up --build
```

## Добавить пользователя

```sh
curl http://localhost:8080/users -X POST -d '{ "first_name": "Ivan", "last_name": "Ivanov" }'
```

## Добавить подписку

С id, полученным из предыдущего запроса, вызвать:

```sh
curl http://localhost:8080/subscriptions -X POST -d \
'{
    "user_id": <USER_ID>,
    "service_name": "Yandex Plus",
    "price": 400,
    "start_date": "02-2025"
}'
```

- при добавлении подписки выполняется проверка даты на то, если она указана в будущем. Например, если текущая дата `07-2025`, а переданное значение `08-2025`, endpoint вернет ошибку

## Вернуть общую стоимость всех подписок в период startDate до endDate

```sh
curl http://localhost:8080/subscriptions/cost?userId=<USER_ID>&serviceName=<SERVICE_NAME>&startDate=date_1&endDate=date_2
```

- `endDate` можно опустить, тогда будет установлена текущая дата
