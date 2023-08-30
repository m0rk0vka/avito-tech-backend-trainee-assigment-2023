FROM golang:1.21.0 as builder
RUN mkdir /app

WORKDIR /app

COPY ./ ./

RUN go mod tidy && \
    go build

FROM ubuntu:22.04 as runner

COPY --from=builder /app/avito-tech-backend-trainee-assigment-2023 ./

ENTRYPOINT ["./avito-tech-backend-trainee-assigment-2023"]
