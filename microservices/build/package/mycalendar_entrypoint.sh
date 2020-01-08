#!/bin/bash
envsubst < config_template.json > config.json
echo "Starting mycalendar"
cat config.json
chmod +x wait
./wait || exit 1
echo "Wait 5 seconds for postgress initialization"
sleep 5
./mycalendar || exit 1
