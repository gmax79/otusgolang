#!/bin/bash
./wait
echo "Starting tests"
./testapp
./testgrpc
./testmq
echo "Test via godog"
./testgodog
