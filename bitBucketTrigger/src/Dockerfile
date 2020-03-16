FROM golang:1.13-alpine as build-go

ENV APP_NAME=gw
RUN apk add -U ca-certificates && adduser -u 1001 scratchuser -D
WORKDIR /go/src/${APP_NAME}

COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gw .

FROM scratch
LABEL MANTAINER "Penny Wise"
ENV APP_NAME=gw

COPY --from=build-go /etc/passwd /etc/passwd
COPY --from=build-go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-go /go/src/${APP_NAME}/gw .

USER scratchuser
CMD ["./gw"]