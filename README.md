# Finch

This is the source code repository for the Finch (https://usefinch.co) project.

## What is it
Finch is a simple service that handles scheduled tasks for your apps, services. You can use the developer friendly API to create your tasks to be completed in a future date, periodically repeated. In plain English, you can create tasks like "Send a request to this URL with this JSON body every morning". That could save you a ton of work when just developing a to-do app.

For now it can send HTTP requests with GET, POST methods and request body of your choice, which can be plain text, form or JSON.

## How to use it
You can use the Swagger page [here](http://swagger.usefinch.co/#/default) to check the endpoints, required fields, expected responses as well as making the actual API calls.

### Getting Access Token
To use the API you will need an access token. You can create yourself one using the `POST /users` endpoint. As a response you will receive the token ([JWT](https://en.wikipedia.org/wiki/JSON_Web_Token)).

### Authenticating API calls
As you can see in the [Swagger](http://swagger.usefinch.co/#/default) page, all endpoints are expecting `Authorization` header. Value of the `Authorization` header must be in the format `Bearer ACCESS_TOKEN`, where `ACCESS_TOKEN` is the one obtained with the step above.

## Running it Yourself
### Using Docker and Docker Compose
You need to have `docker` and `docker-compose` installed on your system. After that it is as simple as running `docker-compose up`. Finch will be listening on `127.0.0.1:8081`, so you can access the API using URLs like `http://localhost:8081/v1/users`.

Finch is using MongoDB and Redis for storing its data, and those services are configured to use volumes in `docker-compose.yml`. So your data will remain in `mongo-data` and `redis-data` folders even if you remove the containers.

## License

[MIT](LICENSE)
