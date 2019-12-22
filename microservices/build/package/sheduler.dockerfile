FROM golang:1.13.5 as builder
COPY . /app
WORKDIR /app/cmd/sheduler
RUN GOOS=linux go build -o sheduler .

FROM ubuntu:18.04
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /root
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait ./wait
COPY build/package/sheduler_entrypoint.sh cmd/sheduler/config_template.json ./ 
COPY --from=builder /app/cmd/sheduler/sheduler ./
ENTRYPOINT "./sheduler_entrypoint.sh"
