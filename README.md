# Effective Mobile test assesment

Порт приложения `APP_PORT=3001` (указан в конфиге .env)

```bash
docker compose up --build
```

## Добавить пользователя

```bash
curl http://localhost:3001/users -X POST -d \
'{
    "first_name": "Ivan",
    "last_name": "Ivanov"
}'
```

## Добавить подписку

С id, полученным из предыдущего запроса, вызвать:

```bash
curl http://localhost:3001/subscriptions -X POST -d \
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
curl http://localhost:3001/subscriptions/cost
```

Со следующими GET-параметрами:

```
userId=<USER_ID>
serviceName=<SERVICE_NAME>
startDate=01-2025
endDate=07-2025
```

- `endDate` можно опустить, тогда будет установлена текущая дата
