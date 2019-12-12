gen:
	cd api/pbcalendar && ./genpb.sh

main:
	cd cmd/mycalendar && go build -o ../../mycalendar

testapp:
	cd cmd/testapp && go build -o ../../testapp

testgrpc:
	cd cmd/testgrpc && go build -o ../../testgrpc

sender:
	cd cmd/sender &&  go build -o ../../sender

sheduler:
	cd cmd/scheduler &&  go build -o ../../scheduler

tests: testapp testgrpc

rabbitmq: sender scheduler

all: gen main testapp testgrpc sender sheduler

