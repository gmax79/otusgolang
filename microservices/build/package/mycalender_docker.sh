#!/bin/sh

envsubs < config_template.json > config.json
./mycalendar
