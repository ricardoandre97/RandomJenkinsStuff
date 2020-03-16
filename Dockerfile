FROM golang:1.13-alpine as build-go

ENV APP_NAME=jobCreator
WORKDIR /go/src/${APP_NAME}

COPY main.go go.mod go.sum ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME .

FROM alpine
LABEL MANTAINER "Penny Wise"
ENV APP_NAME=jobCreator

COPY --from=build-go /go/src/${APP_NAME}/$APP_NAME /usr/bin/$APP_NAME
