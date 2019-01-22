FROM golang:1.11

# Deps
RUN go get github.com/dgrijalva/jwt-go
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq

# Copy app
RUN mkdir /app 
ADD . /app/
WORKDIR /app 

# Build
RUN go build -o main .

# Run
EXPOSE 3000
CMD ["./main"]
