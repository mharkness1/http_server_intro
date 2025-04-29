# Intro to HTTP Servers in Golang

## Project Summary

This project was written alongside the course "Learn HTTP Servers in Go" by boot.dev with project direction given by the course and testing suites provided for each lesson. It is written with the net/http standard library.

The project is a simplified mock of a 'twitter'-eque platform called 'Chirpy' where users are registered and they can post 'chirps'. Currently implemented: user account creation, chirp creation, login with hashed passwords, authentication on posting with JWT and refresh tokens, and authorization (with a fictional premium project 'chirpy red').

Additional tools used: SQLC (ADD LINK) for autogeneration of typed interfaces and dababase methods; Goose (ADD LINK) as a database migration tool.

## Local Set-up
1. Copy repo

``` git clone <repo> <directory> ```

2. Get dependencies

``` go mod download ```

3. Environment variables

4. Database Tools

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