mycalendar:
	cd scripts && ./genpb.sh
	cd cmd/mycalendar && go build -o mycalendar

sender:
	cd cmd/sender &&  go build -o sender

sheduler:
	cd cmd/sheduler &&  go build -o sheduler

testapp:
	cd cmd/tests/testapp && go build -o testapp

testgrpc:
	cd cmd/tests/testgrpc && go build -o testgrpc

testmq:
	cd cmd/tests/testmq && go build -o testmq

testgodog: 
	cd cmd/tests/integration_tests && godog -o testgodog

app: mycalendar sender sheduler

tests: testapp testgrpc testmq testgodog

all: app tests

run: app
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

docker-run-tests:
	set -e; \
	export COMPOSE_IGNORE_ORPHANS=true; \
	exit_code=0; \
	cd deployments; \
	docker-compose -f docker-compose.yml up -d; \
	docker-compose -f docker-compose.tests.yml up || exit_code=$$?; \
	docker-compose -f docker-compose.yml -f docker-compose.tests.yml down; \
	docker rmi gmax079/practice:mycalendar gmax079/practice:sheduler gmax079/practice:sender gmax079/practice:tests
	echo "integration_tests result: $$exit_code"; \
	exit $$exit_code; \

docker-stop:
	cd deployments && docker-compose -f docker-compose.yml -f docker-compose.tests.yml down

# requirement command from hometask
test: docker-run-tests

docker-push:
	docker push gmax079/practice:mycalendar \
	&& docker push gmax079/practice:sheduler \
	&& docker push gmax079/practice:sender \
	&& docker push gmax079/practice:tests

clean: 
	rm -f cmd/mycalendar/mycalendar cmd/sender/sender cmd/sheduler/sheduler
	rm -f cmd/tests/testapp/testapp cmd/tests/testgrpc/testgrpc cmd/tests/testmq/testmq cmd/tests/integration_tests/testgodog

