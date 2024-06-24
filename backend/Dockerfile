FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY .env .env

RUN CGO_ENABLED=0 GOOS=linux go build -o /playlist-maker

EXPOSE 8080

CMD [ "/playlist-maker" ]

