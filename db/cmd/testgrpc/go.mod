module github.com/gmax79/otusgolang/db/cmd/testgrpc

go 1.13

replace github.com/gmax79/otusgolang/db/cmd/mycalendar/pbcalendar => ../mycalendar/pbcalendar

require (
	github.com/gmax79/otusgolang/db/cmd/mycalendar/pbcalendar v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	google.golang.org/grpc v1.25.1
)
