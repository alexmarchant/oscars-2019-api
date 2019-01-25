# Oscars 2019 Api

## Routes

### Users

#### GET /users

Gets all users and their picks

Expected response:

```json
[
  {
    "id": 1,
    "email": "alexjmarchant@gmail.com",
    "picks": {
      "Best Picture": "Spider-Man: Into the Spider-Verse",
      "Best Director": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
    }
  },
  {
    "id": 2,
    "email": "larsonlaidlaw@gmail.com",
    "picks": {
      "Best Picture": "Spider-Man: Into the Spider-Verse",
      "Best Director": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
    }
  }
]
```

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
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoiYWxleGptYXJjaGFudEBnbWFpbC5jb20ifQ.lTKwBXQ09u7JEscdJLDMidHLYLOBvKym8Or7UWsJGXo"
}
```

#### GET /users/current-user

*REQUIRES AUTH*

Get's current user info

Expected response:

```json
{
  "id": 13,
  "email": "alexjmarchant@gmail.com",
  "admin": true
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

### Picks

#### GET /users/current-user/picks

*REQUIRES AUTH*

Gets the current user's picks

Expected response:

```json
{
  "Best Picture": "Spider-Man: Into the Spider-Verse",
  "Best Director": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
}
```

#### POST /users/current-user/picks

*REQUIRES AUTH*

Creates/updates the current user's picks

Expected request:

```json
{
  "Best Picture": "Spider-Man: Into the Spider-Verse",
  "Best Director": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
}
```

Expected response:

Status 200

### Winners

#### POST /winners

*REQUIRES AUTH & ADMIN*

Updates winning picks

Expected request:

```json
{
  "Best Picture": "Spider-Man: Into the Spider-Verse",
  "Best Director": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
}
```

Expected response:

Status 200

## Auth

Must pass token in a header: `Authorization: Bearer <token>`

## Errors

All errors respond with an error message and an appropriate status code:

```json
{
  "error": "Example error message"
}
```
