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

JWT payload:

```json
{
  "admin": true,
  "email": "alexjmarchant@gmail.com"
}
```

### Selection

#### POST /users/:id/selections

Creates new selections

Auth: Must pass token as `Authorization: Bearer <token>`

Expected request:

```json
[
  {
    "category": "Best Picture",
    "selection": "Spider-Man: Into the Spider-Verse"
  },
  {
    "category": "Best Director",
    "selection": "Bob Persichetti, Peter Ramsey, Rodney Rothman"
  }
]
```

Expected response:

Status 201
