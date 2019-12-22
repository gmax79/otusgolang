#!/bin/sh
envsubst < config_template.json > config.json
echo "Starting mycalendar"
cat config.json
./mycalendar
