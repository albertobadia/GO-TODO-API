# TODO API WITH GO

## Description

This is a simple API with Go, the goal is to manage a TODO list.

## Starting 🚀

### Pre-requisites 📋

- Docker
- Docker-compose
- OpenSSL

### Installation 🔧

1. Clone the repository
2. Generate a secret key with the following command:

```
openssl rand -base64 32
```

3. Create a file called `.env` in the root of the project with the following content:

```
SECRET_KEY={YOUR_SECRET_KEY}
POSTGRES_USER={YOUR_POSTGRES_USER} Optional (default: todo)
POSTGRES_PASSWORD={YOUR_POSTGRES_PASSWORD} Optional (default: todo_pass)
POSTGRES_DB={YOUR_POSTGRES_DB} Optional (default: todos)
POSTGRES_HOST={YOUR_POSTGRES_HOST} Optional (default: postgres)
POSTGRES_PORT={YOUR_POSTGRES_PORT} Optional (default: 5432)
```

4. Run the following command:

```
docker-compose up -d --build
```

5. Once the database is running, the api will be available at <http://localhost:8080>

## Testing ⚙️

### Pre-requisites 📋

- Go installed
- Install the dependencies with the following command:
```
go mod download
```
- Set IS_TESTING environment variable to true with the following command:
```
export IS_TESTING=true
```

### Run the tests

```
go test ./tests -v
```

## Swagger

The API documentation is available at <http://localhost:8080/swagger-ui/>

## Endpoints

### Register

```
POST /register HTTP/1.1 Host: localhost:8080 Content-Type: application/json

{
    "username": "test",
    "password": "test"
}
```

### Login

```
POST /login HTTP/1.1 Host: localhost:8080 Content-Type: application/json

{
    "username": "test",
    "password": "test"
}
```

### Create a TODO

```
POST /todos HTTP/1.1 Host: localhost:8080 Content-Type: application/json Authorization: Bearer {TOKEN}
{"title": "Test"}

Response:
{
    "id": 3f1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
    "user_id": 3r1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
    "title": "Test",
    "is_done": false
}
```

### Get all TODOs

```
GET /todos HTTP/1.1 Host: localhost:8080 Content-Type: application/json Authorization: Bearer {TOKEN}
Response:
[
    {
        "id": 3f1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
        "user_id": 3r1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
        "title": "Test",
        "is_done": false
    }
]
```

### Get a TODO

```
GET /todos/{id} HTTP/1.1 Host: localhost:8080 Content-Type: application/json Authorization: Bearer {TOKEN}
Response:
{
    "id": 3f1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
    "user_id": 3r1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
    "title": "Test",
    "is_done": false
}
```

### Update a TODO

```
PUT /todos/{id} HTTP/1.1 Host: localhost:8080 Content-Type: application/json Authorization: Bearer {TOKEN}
{"title": "Test", "description": "Test"}

Response:
{
    "id": 3f1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
    "user_id": 3r1b1c1a-1b1c-4f1b-1c1a-1b1c1f1b1c1a,
    "title": "Test",
    "is_done": false
}
```

### Delete a TODO

```
DELETE /todos/{id} HTTP/1.1 Host: localhost:8080 Content-Type: application/json Authorization: Bearer {TOKEN}
Response: 204 No Content
```

## Built with 🛠️

- [Go](https://golang.org/) - Programming language
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Docker](https://www.docker.com/) - Containerization platform
- [Docker-compose](https://docs.docker.com/compose/) - Tool for defining and running multi-container Docker applications
- Other Go packages like mux, bcrypt, etc
