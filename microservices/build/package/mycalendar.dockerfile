FROM golang:1.13.5 as builder
COPY . /app
WORKDIR /app/cmd/mycalendar
RUN GOOS=linux go build -o mycalendar .

FROM ubuntu:18.04
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /root
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait ./wait
COPY build/package/mycalendar_entrypoint.sh cmd/mycalendar/config_template.json ./ 
COPY --from=builder /app/cmd/mycalendar/mycalendar ./
ENTRYPOINT "./mycalendar_entrypoint.sh"
