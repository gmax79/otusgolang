FROM golang:1.13.5 as builder
RUN go get github.com/DATA-DOG/godog/cmd/godog
COPY . /app
WORKDIR /app/cmd/tests/testapp
RUN GOOS=linux go build -o testapp .
WORKDIR /app/cmd/tests/testgrpc
RUN GOOS=linux go build -o testgrpc .
WORKDIR /app/cmd/tests/testmq
RUN GOOS=linux go build -o testmq .
WORKDIR /app/cmd/tests/integration_tests
RUN GOOS=linux godog -o testgodog

FROM ubuntu:18.04
WORKDIR /root
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.2.1/wait ./wait
COPY build/package/tests_entrypoint.sh ./
COPY --from=builder /app/cmd/tests/testapp/testapp ./
COPY --from=builder /app/cmd/tests/testgrpc/testgrpc ./
COPY --from=builder /app/cmd/tests/testmq/testmq ./
COPY --from=builder /app/cmd/tests/integration_tests/features ./features/
COPY --from=builder /app/cmd/tests/integration_tests/testgodog ./
ENTRYPOINT "./tests_entrypoint.sh"
