#!/bin/sh

cd ../api/pbcalendar
protoc -I. -I/usr/local/include --go_out=plugins=grpc:. mycalendar.proto
