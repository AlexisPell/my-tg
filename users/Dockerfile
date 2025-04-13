FROM golang:latest
# RUN go get -u github.com/cosmtrek/air
RUN go install github.com/air-verse/air@latest
WORKDIR /app
ENTRYPOINT ["air"]