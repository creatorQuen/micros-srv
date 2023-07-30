FROM golang:1.20-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerAppication ./cmd/api

RUN chmod +x /app/brokerAppication


FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerAppication /app

CMD [ "/app/brokerAppication" ]