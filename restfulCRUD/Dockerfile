
FROM golang:latest AS builder

ENV GIN_MODE=release

WORKDIR /app

COPY . .

RUN go build -o main cmd/app/main.go


# Run stage
FROM alpine:3.16

COPY --from=builder /app/main ./

EXPOSE 8080

CMD [ "/main" ]