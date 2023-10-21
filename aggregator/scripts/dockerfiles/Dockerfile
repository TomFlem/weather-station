FROM golang:1.21 AS builder
RUN apt-get update && apt-get install -y make gcc
WORKDIR /aggregator
COPY . .
RUN make
FROM debian:stable-slim
LABEL maintainer="tom"
COPY --from=builder /aggregator/weather-aggregator /weather-aggregator
ENTRYPOINT ["./weather-aggregator"]