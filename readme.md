# Go API example with JWT auth 

### Boiler plate code for APIs protected by username/login and javascript web token with a 10 minute expiry time.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/feliperyan/golang_api_boilerplate) [![Go Report Card](https://goreportcard.com/badge/github.com/heroku-examples/golang-jwt-api-boilerplate)](https://goreportcard.com/badge/github.com/heroku-examples/golang-jwt-api-boilerplate)

**Companion app, which provides the UI for this API server**
- https://github.com/feliperyan/react_materialui_boilerplate

**Features**
- New users get saved to Postgres.
- Set Env Var `NEEDS_AUTH=yes` on Heroku for auth
- If you want auth also set `TOKEN_PASSWORD=yourlongsecretkeyhere` so you can generate tokens

**Endpoints**:

1. api/user/login
2. api/user/new
3. api/quote

The `login` and `new` take a JSON payload of `{"email": "valar@morghulis.com", "password": "valardohaeris"}`

While `quote` requires an Authorisation Header with the Bearer token issued by user creation or log in (if NEEDS_AUTH == yes)

>This repo is heavily based on Adigun Olalekan's excellent post here:
>https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

again...
