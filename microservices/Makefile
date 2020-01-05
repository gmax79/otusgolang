mycalendar:
	cd api/pbcalendar && ./genpb.sh
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

testgodog: 
	cd cmd/integration_tests && godog -o testgodog

tests: testapp testgrpc testmq testgodog

all: mycalendar sender sheduler tests

run: 
	cd cmd/mycalendar && ./mycalendar &
	cd cmd/sender && ./sender &
	cd cmd/sheduler && ./sheduler &

run-tests: tests
	cd cmd/testapp && ./testapp
	cd cmd/testgrpc && ./testgrpc
	cd cmd/testmq && ./testmq
	cd cmd/integration_tests && ./testgodog

stop:
	build/kill.sh sheduler
	build/kill.sh sender
	build/kill.sh mycalendar

docker-mycalendar:
	docker build -f build/package/mycalendar.dockerfile -t mycalendar . \
	&& docker tag mycalendar gmax079/practice:mycalendar

docker-sheduler:
	docker build -f build/package/sheduler.dockerfile -t sheduler . \
	&& docker tag sheduler gmax079/practice:sheduler

docker-sender:
	docker build -f build/package/sender.dockerfile -t sender . \
	&& docker tag sender gmax079/practice:sender

docker-tests:
	docker build -f build/package/tests.dockerfile -t tests . \
	&& docker tag tests gmax079/practice:tests

docker-all: docker-mycalendar docker-sheduler docker-sender docker-tests

docker-run:
	cd deployments && docker-compose up

docker-stop:
	cd deployments && docker-compose down

docker-push:
	docker push gmax079/practice:mycalendar \
	&& docker push gmax079/practice:sheduler \
	&& docker push gmax079/practice:sender \
	&& docker push gmax079/practice:tests

clean: 
	-rm cmd/mycalendar/mycalendar cmd/testapp/testapp cmd/testgrpc/testgrpc cmd/testmq/testmq cmd/integration_tests/testgodog cmd/sender/sender cmd/sheduler/sheduler 2>/dev/null 
