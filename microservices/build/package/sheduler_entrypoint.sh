#!/bin/bash
envsubst < config_template.json > config.json
echo "Starting sheduler"
cat config.json
chmod +x wait
./wait || exit 1
./sheduler || exit 1
