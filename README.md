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
    "picks": [
      {
        "category": "Best Picture",
        "pick": "Spider-Man: Into the Spider-Verse"
      },
      {
        "category": "Best Director",
        "pick": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
      }
    ]
  },
  {
    "id": 2,
    "email": "larsonlaidlaw@gmail.com",
    "picks": [
      {
        "category": "Best Picture",
        "pick": "Spider-Man: Into the Spider-Verse"
      },
      {
        "category": "Best Director",
        "pick": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
      }
    ]
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

### Picks

#### GET /users/:id/picks

Gets a single user's picks

Expected response:

```json
[
  {
    "category": "Best Picture",
    "pick": "Spider-Man: Into the Spider-Verse"
  },
  {
    "category": "Best Director",
    "pick": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
  }
]
```

#### PUT /users/:id/picks

*REQUIRES AUTH*

Creates/updates a user's picks

Expected request:

```json
[
  {
    "category": "Best Picture",
    "pick": "Spider-Man: Into the Spider-Verse"
  },
  {
    "category": "Best Director",
    "pick": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
  }
]
```

Expected response:

Status 200

## Auth

Must pass token in a header: `Authorization: Bearer <token>`

JWT payload:

```json
{
  "id": 1,
  "email": "alexjmarchant@gmail.com",
  "admin": true
}
```

## Errors

All errors respond with an error message and an appropriate status code:

```json
{
  "error": "Example error message"
}
```
