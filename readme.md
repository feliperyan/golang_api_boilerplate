# Go API example with JWT auth

### This repo is heavily based on Adigun Olalekan's excellent post here: https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

Boiler plate code for APIs protected by username/login and javascript web token with a 10 minute expiry time.

New users get saved to Postgres.

Endpoints:

1. api/user/login
2. api/user/new
3. api/dummy

The `login` and `new` take a JSON payload of `{"email": "valar@morghulis.com", "password": "valardohaeris"}`

While `dummy` required an Authorisation Header with the Bearer token issues by user creation or log in.

