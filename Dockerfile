FROM golang:1.12-stretch

WORKDIR /src/

COPY ./ /src/

RUN go build -o bin/bailer

FROM debian:stretch-slim

RUN apt-get update && \
  apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*

COPY --from=0 /src/bin/bailer /usr/local/bin/bailer

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/bailer"]
