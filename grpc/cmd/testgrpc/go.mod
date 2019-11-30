module github.com/gmax79/otusgolang/grpc/cmd/testgrpc

go 1.13

replace github.com/gmax79/otusgolang/grpc/cmd/mycalendar/pbcalendar => ../mycalendar/pbcalendar

require (
	github.com/gmax79/otusgolang/grpc/cmd/mycalendar/pbcalendar v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.25.1
)
