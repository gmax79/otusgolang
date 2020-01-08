#!/bin/bash
chmod +x wait
./wait || exit 1
echo "Wait 5 seconds before starting"
sleep 5
echo "Starting tests"
./testapp || exit 1
./testgrpc || exit 1
./testmq || exit 1
echo "Tests via godog"
./testgodog || exit 1
