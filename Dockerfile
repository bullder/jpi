FROM golang:latest AS builder

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

FROM alpine:latest as app

RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
COPY ./docs ./
COPY ./public ./public

RUN chmod +x ./main

EXPOSE 8080

CMD ./main
