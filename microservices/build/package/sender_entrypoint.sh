#!/bin/bash
envsubst < config_template.json > config.json
echo "Starting sender"
cat config.json
chmod +x wait
./wait && ./sender
