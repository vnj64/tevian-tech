# Работа с проектом.
### Запуск проекта
### `docker compose up -d`

### Остановка проекта
### `docker compose down (-v optional)`

# Настройка переменных окружения.
В файле .env необходимо поменять переменные **_CLOUD_LOGIN_** и **_CLOUD_PASSWORD_** на свои данные (находящиеся в .env также работают). 

# Описание возможных настроек сервиса.
Все, что вы можете конфигурировать находится в .env, вся остальная работа происходит в фоновом режиме, а именно:
- Получение **_access_token_** для работы с методами FaceCloud. Пользователю не нужно беспокоиться о прохождении авторизации на постороннем сервисе.
- Обращение к запросу /api/v1/detect находящемуся в API FaceCloud.

Также можно редактировать следующие данные:
- Логин от FaceCloud сервиса
- Пароль от FaceCloud сервиса
- Пароль базы данных
- Юзернейм базы данных
- Хост базы данных
- Порт базы данных
- Название базы данных
- Базовый адрес обращения к API FaceCloud

# Описание API сервиса.
* Добавление задания
```http request
POST /api/v1/task/
```
* Добавление изображения в задание
```http request
POST /api/v1/task/:id/upload_image
```
* Запуск обработки задания
```http request
POST /api/v1/task/:id/start_task
```
* Получение задания
```http request
GET /api/v1/task/:id
```
* Удаление задания
```http request
DELETE /api/v1/task/:id
```

# Комментарии к проекту
В этом проекте можно было бы использовать более простую архитектуру проекта, но заглядывая в перспективу расширения такого API сервиса у нас зачастую будет возникать необходимость в расширении его функционала, поэтому выбор был сделан в пользу DDD. А теперь к сути!

На этом умные слова об этом проекте заканчиваются. Надеюсь увидеться на собеседовании!)