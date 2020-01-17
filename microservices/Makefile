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

devenv-up:
	docker-compose -f deployments/docker-compose.dev.yml up -d

devenv-down:
	docker-compose -f deployments/docker-compose.dev.yml down

run: app devenv-up
	sleep 10 #wait rabbitmq and postgres
	cd cmd/mycalendar && ./mycalendar &
	cd cmd/sender && ./sender &
	cd cmd/sheduler && ./sheduler &

run-tests: tests # required make run before
	-cd cmd/tests/testapp && ./testapp
	-cd cmd/tests/testgrpc && ./testgrpc
	-cd cmd/tests/testmq && ./testmq
	-cd cmd/tests/integration_tests && ./testgodog
	make clean-tests

stop:
	-pkill -f sheduler
	-pkill -f sender
	-pkill -f mycalendar
	make devenv-down
	make clean-app

devenv-run:
	cd deployments && docker-compose -f docker-compose.dev.yml up

devenv-stop:
	cd deployments && docker-compose -f docker-compose.dev.yml down

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
	cp /etc/localtime /tmp/localtime
	cd deployments && docker-compose -f docker-compose.yml up

docker-stop:
	cd deployments && docker-compose -f docker-compose.yml -f docker-compose.tests.yml down

docker-run-tests:
	set -e; \
	export COMPOSE_IGNORE_ORPHANS=true; \
	cp /etc/localtime /tmp/localtime; \
	cd deployments && docker-compose -f docker-compose.tests.yml up

docker-cicd-tests:
	set -e; \
	export COMPOSE_IGNORE_ORPHANS=true; \
	cp /etc/localtime /tmp/localtime; \
	exit_code=0; \
	docker-compose -f deployments/docker-compose.yml up -d; \
	docker-compose -f deployments/docker-compose.tests.yml up || exit_code=$$?; \
	docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.tests.yml down; \
	make docker-clean; \
	echo "integration_tests result: $$exit_code"; \
	exit $$exit_code;

docker-prom-tests:
	make docker-run-tests
	make docker-run-tests
	make docker-run-tests

# requirement command from hometask
test: docker-cicd-tests

docker-push:
	docker push gmax079/practice:mycalendar \
	&& docker push gmax079/practice:sheduler \
	&& docker push gmax079/practice:sender \
	&& docker push gmax079/practice:tests

docker-clean:
	docker rmi -f gmax079/practice:mycalendar gmax079/practice:sheduler gmax079/practice:sender gmax079/practice:tests

clean-app:
	rm -f cmd/mycalendar/mycalendar cmd/sender/sender cmd/sheduler/sheduler

clean-tests:
	rm -f cmd/tests/testapp/testapp cmd/tests/testgrpc/testgrpc cmd/tests/testmq/testmq cmd/tests/integration_tests/testgodog

clean: clean-app clean-tests
