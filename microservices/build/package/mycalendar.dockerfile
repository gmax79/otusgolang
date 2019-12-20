FROM golang:1.13.5 as builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/mycalendar
RUN GOOS=linux go build -o mycalendar .

FROM alpine:latest  
#RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /app/cmd/mycalendar .
CMD ["./mycalendar"] 
