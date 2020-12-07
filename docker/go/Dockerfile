FROM golang:1.14 as build

ARG APP_NAME
ENV APP_NAME "${APP_NAME}"

WORKDIR /go/src/app

COPY . .

RUN go build -o bin/${APP_NAME} ./cmd/${APP_NAME}

CMD bin/${APP_NAME}