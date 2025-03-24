# FROM golang:1.24 AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go build -o main .

# RUN go install github.com/cosmtrek/air@latest

# CMD ["./main"]


FROM golang:latest
# RUN go get -u github.com/cosmtrek/air
RUN go install github.com/air-verse/air@latest
WORKDIR /app
ENTRYPOINT ["air"]