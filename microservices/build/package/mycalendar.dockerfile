FROM golang:1.13.5 as builder
COPY . /app
WORKDIR /app/cmd/mycalendar
RUN GOOS=linux go build -o mycalendar .

FROM alpine:latest  
RUN \
  apt-get update \
  && apt-get -y install gettext-base \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

#RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/cmd/mycalendar .
COPY mycalendar/config_template.json .
COPY mycalendar_docker.sh .
CMD ["mycalendar_docker.sh"]
