# Чекин Игорь, Сервис-ориентированные архитектуры, Домашнее задание 7

## Общая структура проекта

В папке `MafiaGQL_server` содержатся все основные файлы приложения. Обозначим назначения пакетов внутри модуля (название пакетов совпадают с названиями папок).
* `config.go` - здесь задаются константы, используемые остальными пакетами модуля
* `db` - имплементирует интерфейс для взаимодействия с базой данных (в моем случае с `MongoDB`)
* `graph` - содержит описание структур данных GraphQL, а также обработчики запросов
* `utils` - общие вспомогательные функции


В качестве базы данных используется NoSQL-решение `MongoDB`, в связи с простотой интерфейса и нативностью сериализации данных из `json` в Mongo и обратно.

Для локального запуска сервера необходимо наличие `go` версии `1.17 - 1.18`. Требуется перейти в папку `MafiaGQL_server` и ввести команду `go run .`
[Ссылка на докер-образ](https://hub.docker.com/layers/soa-images/alucardik/soa-images/MafiaREST/images/sha256-3eedf991a93496601bce41eb5d14c83a9c94ee3b3609ab027532f9c2e2ab9b4a?context=explore) (`alucardik/soa-images:MafiaGQL-server`)

> ВНИМАНИЕ: для работы сервера нужна поднятая база данных, это может усложнять локальный запуск

Тем не менее, настоятельно рекомендуется запускать все сервисы через `docker-compose`, так как между ними уже настроены все связи. Достаточно выполнить следующие команды из корневой папки:

```bash
docker-compose up --build
```

Далее можно отправлять запросы на `localhost:8080/query`.

## Описание структур данных

`GameSession` - представляет сессию игрока. Представляется следующим `JSON`:
```
{
    "_id": {идентификатор сессии}, // строка
    "name": {название сессии}, // строка
    "participants": {участники сесси} // список строк
    "comments": {комментарии к сессии} // список структур Comment
}
```

`Comment` - представляет комментарий игрока к сессии. Представляется следующим `JSON`:
```
{
	"sessionId": {идентификатор сессии} // строка
	"author": {ник автора комментария} // строка
	"contents": {содержимое комментария} // строка
}
```

## Описание API сервиса

### Запросы
Сервис поддерживает единый запрос `sessions` со следующими параметрами:

```sessions(ongoing: Boolean = true, sessionId: ID): [GameSession!]```, где

* `ongoing` указывает искать ли активные сессии (`true`, значение по умолчанию) или неактивные (`false`)
* при указании `sessionId` выводится информация о конкретной сессии, параметр `ongoing` в этом случае игнорируется 

Примеры `GraphQL` запросов:

```graphql
# вывести _id, названия и комментарии к активным сессиям
query {
    sessions {
        _id
        name
        comments {
            author
            contents
        }
    }
}

# аналогичный запрос к прошедшим сессиям
query {
    sessions(ongoing: false) {
        _id
        name
        comments {
            author
            contents
        }
    }
}

# аналогичный запрос к сессии с конкретным id
query {
    sessions(sessionId: "628a0bc94a928ca0f0d4500f") {
        _id
        name
        comments {
            author
            contents
        }
    }
}
```

### Мутации
Сервис поддерживает следующие мутации:

```graphql
# создать новую сессию
startSession(input: NewGameSession!): GameSession!

input NewGameSession {
    name: String!
    initiator: String!
}
# -----------------------------------------------

# добавить участника к активной сессии
addParticipant(input: NewParticipant!): GameSession

input NewParticipant {
    sessionId: String!
    userId: String!
}
# -----------------------------------------------

# добавить комментарий к сессии  
addComment(input: NewComment!): String!

input NewComment {
    sessionId: String!
    author: String!
    contents: String!
}
# -----------------------------------------------
    
# завершить сессию, после этого нельзя будет добавлять участников
endSession(sessionId: String!): String!
```

Примеры `GraphQL` мутаций:

```graphql
# создать активную сессию "MySessionName" с единственным игроком "MyPlayerName"
mutation {
    startSession(input: {name: "MySessionName", initiator: "MyPlayerName"}) {
        _id
        name
    }
}

# добавить участника "playerName" к сессии с "sessionId"
mutation {
    addParticipant(input: {sessionId: "628a0bc94a928ca0f0d4500f", userId: "playerName"}) {
        name
        participants
    }
}

# добавить комментарий к сессии, в ответ приходит уведомление о статусе операции
mutation {
    addComment(input: {sessionId: "628a0bc94a928ca0f0d4500f", author: "SomePlayer", contents: "My very first comment" })
}

# завершить сессию со следующим id, в ответ приходит уведомление о статусе операции
mutation {
    endSession(sessionId: "628a0bc94a928ca0f0d4500f")
}
```
