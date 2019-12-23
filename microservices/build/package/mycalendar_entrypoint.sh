#!/bin/bash
envsubst < config_template.json > config.json
echo "Starting mycalendar"
cat config.json
chmod +x wait
./wait && ./mycalendar
