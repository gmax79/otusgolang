gen:
	cd api/pbcalendar && ./genpb.sh

main: gen
	cd cmd/mycalendar && go build -o mycalendar

testapp:
	cd cmd/testapp && go build -o testapp

testgrpc:
	cd cmd/testgrpc && go build -o testgrpc

sender:
	cd cmd/sender &&  go build -o sender

sheduler:
	cd cmd/sheduler &&  go build -o sheduler

testmq:
	cd cmd/testmq && go build -o testmq

tests: testapp testgrpc testmq

rabbitmq: sender scheduler

all: main tests sender sheduler

clean: 
	-rm cmd/mycalendar/mycalendar cmd/testapp/testapp cmd/testgrpc/testgrpc cmd/testmq/testmq cmd/sender/sender cmd/sheduler/sheduler 2>/dev/null 
