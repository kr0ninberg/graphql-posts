# Учебный проект на Go с использованием GraphQL

GraphQL-сервис на Go для управления постами и комментариями


## Функционал

### Посты
- Просмотр списка постов
- Получение поста по ID с комментариями
- Автор поста может запретить комментарии

### Комментарии
- Иерархические (вложенные) комментарии без ограничения глубины
- Пагинация комментариев и ответов (`limit`, `offset`)
- Ограничение на длину текста комментария: до 2000 символов


## Архитектура

- Backend: Go
- GraphQL: gqlgen
- Хранилище: 
  - `In-Memory`
  - `PostgreSQL` (переключается через конфиг config.go)
- Docker: `Docker + docker-compose`


## Запуск

```
git clone https://github.com/kr0ninberg/graphql-posts.git
cd graphql-posts
# .env file with dsn link required
# you can use "mv .env.example .env" or provide your own dsn 
docker-compose up
# use "docker-compose up qraphql-service" if you dont want to 
# run provided postgres container
```
