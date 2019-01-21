# Oscars 2019 Api

## Routes

### Users

#### POST /users

Creates a new user

Expected request:

```json
{
  "email": "alexjmarchant@gmail.com",
  "password": "12345678",
  "passwordConfirmation": "12345678"
}
```

Expected response:

```json
{
  "id": 10
}
```

### Auth

#### POST /tokens

Creates a new JWT token

Expected request:

```json
{
  "email": "alexjmarchant@gmail.com",
  "password": "12345678"
}
```

Expected response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoiYWxleGptYXJjaGFudEBnbWFpbC5jb20ifQ.lTKwBXQ09u7JEscdJLDMidHLYLOBvKym8Or7UWsJGXo"
}
```
