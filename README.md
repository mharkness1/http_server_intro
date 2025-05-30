# Intro to HTTP Servers in Golang

## Project Summary

This project was written alongside the course "Learn HTTP Servers in Go" by boot.dev with project direction given by the course and testing suites provided for each lesson. It is written with the net/http standard library.

The project is a simplified mock of a 'twitter'-eque platform called 'Chirpy' where users are registered and they can post 'chirps'. Currently implemented: user account creation, chirp creation, login with hashed passwords, authentication on posting with JWT and refresh tokens, and authorization (with a fictional premium project 'chirpy red').

Additional tools used: [SQLC](https://sqlc.dev/) for autogeneration of typed interfaces and dababase methods; [Goose](https://github.com/pressly/goose) as a database migration tool.

The main purpose of this project was to build a web server with additional database functionality, it does include a 'front-end' which is a static page found at localhost:8080 as well as a webpage found at the endpoint /api/metrics which returns the number of times the page has been visited.

## Local Set-up
1. Copy repo

` git clone <repo> <directory> `

2. Get dependencies

You can download the listed dependencies using:

` go mod download `

However, this won't install all necessary dependencies for this project. In particular, the database management tool 'goose' is needed which can be installed with:

`go install github.com/pressly/goose/v3/cmd/goose@latest`

3. Local Database Set-up

This project was developed for use with postgres (version 15). Once installed locally create a database called 'chirpy'.

You can check if you have psql (the defauly cli for postgres) by typing `psql --version` in your terminal. Once installed start a postgres server in the background, enter the psql shell and create a new database using: `CREATE DATABASE chirpy;`

You can check that it has correctly started by entering the database and querying: `SELECT version();`

4. Environment variables

Create a file named: '.env' add to it the following environment variables that are imported and used through the code:

```
DB_URL = postgres://<username>:@localhost:<port>/chirpy?sslmode=disable
PLATFORM = "dev"
SECRET = <Any random string>
POLKA_KEY = "f271c81ff7084ee5b99a5091b42d486e"
```

A random string can easily be generated in terminal using the following command:
`openssl rand -base64 <n>`
where `<n>` is the length of the random string desired e.g, `openssl rand -base64 64`

The Polka Key is the api key of a fictional third party service, to test the authorization is correctly implemented this can be set to any key as long as the authorization header in the corresponding request matches.

PLATFORM allows for access to the /api/reset endpoint, which allowed for boot.dev's testing suite to have a clean database for each submission.

5. Running
- Ensure postgres is running locally on the machine and the local set-up has been completed.
- Using goose bring the database up to latest version using:
``goose postgres postgres://<username>:@localhost:<database port>/chirpy up``
- Build and run the project:
``go build -o out && ./out``
- Test it out with API to the endpoints specified below.

## Available Resources
### Summary
| Endpoint | Method | Description |
| ----------- | ----------- | ----------- |
| /api/chirps | GET | Returns all chirps with query support |
| /api/chirps | POST | Creates, cleans, and validates a chirp and adds it to the database |
| /api/chirps/{chirpID} | GET | Returns a chirp matching the UUID provided |
| /api/chirps/{chirpID} | DELETE | Deletes chirp matching the UUID where the user is the original creator |
| /api/healthz | GET | Returns Status Code 200 and 'OK' if the server is running |
| /api/login | POST | Checks, authenticates and returns a refresh token when a user logs in |
| /api/metrics | GET | Returns HTML with dynamic # of visits counter |
| /api/polka/webhooks | POST | Third party user upgrade mock-up |
| /api/refresh | POST | Returns a refresh token with new expiry |
| /api/reset | POST | If PLATFORM = "dev" resets # of visits counter and deletes all users |
| /api/revoke | POST | Revokes an existing refresh token, without authentication currently |
| /api/users | POST | Creates a user and adds it to the database |
| /api/users | PUT | Authenticates request and updates user information in the database|

### Detailed Descriptions by Endpoint

All requests must contain the relevant authorizations – in particular the relevant JWT (which contains necessary user ID data amongst other properties).

#### 1. /api/chirps - GET & POST
##### GET
This endpoint returns chirps from the database. By default it returns them all ordered by the date created in ascending order (oldest first).

- Queries
This endpoint accepts two URL query parameters:
1. `?author_id=<author id>` where the author ID is the UUID of a specific user that created the chirp.
2. `?sort=<order>` if 'order' is 'desc' then it will return either all or the author's chirps in descending order (newest first). Anything other than 'desc' will return chirps in default order.

- Response Structure

```
[
  {
    "id": 1,
    "body": "<chirp 1>",
    "author_id": <user id of the chirp poster>,
    "created_at": "<created time of chirp 1>",
    "updated_at": "<updated time of chirp 1>"
  },
  {
    "id": 2,
    "body": "<chirp 2>",
    "author_id": <user id of the chirp poster>,
    "created_at": "<created time of chirp 2>",
    "updated_at": "<updated time of chirp 2>"
  }
]
```

##### POST
This endpoint is used to create a new 'chirp'.

- Request Structure

```
{
    "body":"<chirp content>"
}
```

- Response Structure

```
  {
    "id": <chirp id>,
    "body": "<chirp content>",
    "author_id": <user id of the chirp poster>,
    "created_at": "<created time of chirp>",
    "updated_at": "<updated time of chirp>"
  }
```

#### 2. /api/chirps/{chirpID} - GET & DELETE
##### GET
This endpoint simply returns the chirp with the ID matching the provided path.

- Response Structure

```
  {
    "id": <chirp id>,
    "body": "<chirp content>",
    "author_id": <user id of the chirp poster>,
    "created_at": "<created time of chirp>",
    "updated_at": "<updated time of chirp>"
  }
```

##### DELETE

Assuming the necessary authentication it deletes the chirp with the ID provided and returns a status code.

#### 3. /api/healthz - GET

Returns a simple status code.

#### 4. /api/login - POST

Send email and password to which the server responds with the necessary authorizations/tokens as well as user details from the system.

- Request Structure

```
{
    "email":"user email",
    "password":"user password"
}
```

- Response Structure

```
{
    "id":"<user id>",
    "created_at":"<date the user was created>",
    "updated_at":"<date the user details were last updated>",
    "email":"<user email>",
    "token":"<authorization token>",
    "refresh_token":"<refresh token>",
    "is_chirpy_red":"<true/false reflecting premium status>"
}
```

#### 5. /api/metrics - GET

Returns a webpage with the total number of page visits.

#### 6. /api/polka/webhooks - POST

Dummy implementation of a premium service managed by a third party 'polka'. If the provided event is anything other than 'user.upgraded' this will return an error. The .env polka key must also match the provided token in the authorization header of the request.

- Request Structure

```
{
    "event":"user.upgraded",
    "data": {
        "user_id":"<user id to be upgraded>",
    }
}
```

#### 7. /api/refresh - POST

If a refresh token exists and is valid this endpoint will return a fresh token.

- Response Structure

```
{
    "token":"<new refresh token>"
}
```

#### 8. /api/revoke - POST

This end point revokes the refresh token in the header of the request and returns a status code.

#### 9. /api/users - POST & PUT
##### POST
This creates a user and adds them to the database. It simply takes the email and password, and responds with the user details as registered in the database.

- Request Structure

```
{
    "email":"user email",
    "password":"user password"
}
```

- Response Structure

```
{
    "id":"<user id>",
    "created_at":"<date the user was created>",
    "updated_at":"<date the user details were last updated>",
    "email":"<user email>",
    "token":"<authorization token>",
    "refresh_token":"<refresh token>",
    "is_chirpy_red":"<true/false reflecting premium status>"
}
```

##### PUT
This end point updates user info: email and password, after validating using the user id stored in the authorization header of the request to match to the relevant user.

- Request Structure

```
{
    "email":"new user email",
    "password":"new user password"
}
```

- Response Structure

```
{
    "id":"<user id>",
    "created_at":"<date the user was created>",
    "updated_at":"<date the user details were last updated>",
    "email":"<new user email>",
    "token":"<authorization token>",
    "refresh_token":"<refresh token>",
    "is_chirpy_red":"<true/false reflecting premium status>"
}
```