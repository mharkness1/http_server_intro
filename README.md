# Intro to HTTP Servers in Golang

## Project Summary

This project was written alongside the course "Learn HTTP Servers in Go" by boot.dev with project direction given by the course and testing suites provided for each lesson. It is written with the net/http standard library.

The project is a simplified mock of a 'twitter'-eque platform called 'Chirpy' where users are registered and they can post 'chirps'. Currently implemented: user account creation, chirp creation, login with hashed passwords, authentication on posting with JWT and refresh tokens, and authorization (with a fictional premium project 'chirpy red').

Additional tools used: [SQLC](https://sqlc.dev/) for autogeneration of typed interfaces and dababase methods; [Goose](https://github.com/pressly/goose) as a database migration tool.

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

5. Running

## Available Resources

PLACE HOLDERS
mux.HandleFunc("GET /api/healthz", healthzHandler)
mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)
mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)
mux.HandleFunc("POST /api/login", apiCfg.loginHandler)
mux.HandleFunc("POST /api/refresh", apiCfg.refreshTokenHandler)
mux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshTokenHandler)
mux.HandleFunc("PUT /api/users", apiCfg.updateUserInfoHandler)
mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpHandler)
mux.HandleFunc("POST /api/polka/webhooks", apiCfg.upgradeUserHandler)

## Data Formats