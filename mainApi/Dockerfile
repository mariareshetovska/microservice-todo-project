FROM golang:1.18.2-alpine3.14 as build
WORKDIR /go/src
COPY . .
RUN go build -o /app

FROM alpine:3.14
WORKDIR /usr/src
RUN addgroup -S gouser && adduser -S -G gouser gouser
USER gouser:gouser
COPY --from=build /app app
COPY ./migration ./migration

EXPOSE 8080
CMD ["./app" ]