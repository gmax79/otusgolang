gen:
	cd api/pbcalendar && ./genpb.sh

mycalendar: gen
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

all: mycalendar tests sender sheduler

docker-mycalendar: mycalendar
	docker build -f build/package/mycalendar.dockerfile -t mycalendar . \
	&& docker tag mycalendar gmax079/practice:mycalendar

docker-sheduler: sheduler
	docker build -f build/package/sheduler.dockerfile -t sheduler . \
	&& docker tag sheduler gmax079/practice:sheduler

docker-sender: sender
	docker build -f build/package/sender.dockerfile -t sender . \
	&& docker tag sender gmax079/practice:sender

docker-all: docker-mycalendar docker-sheduler docker-sender

docker-push:
	docker push gmax079/practice:mycalendar \
	&& docker push gmax079/practice:sheduler \
	&& docker push gmax079/practice:sender

clean: 
	-rm cmd/mycalendar/mycalendar cmd/testapp/testapp cmd/testgrpc/testgrpc cmd/testmq/testmq cmd/sender/sender cmd/sheduler/sheduler 2>/dev/null 
