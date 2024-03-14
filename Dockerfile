FROM golang:alpine3.19 AS build

WORKDIR /app

COPY . .

RUN go build -o /regex cmd/main.go

CMD [ "/regex" ]
