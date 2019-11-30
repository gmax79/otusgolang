module github.com/gmax79/otusgolang/grpc/cmd/mycalendar

replace github.com/gmax79/otusgolang/grpc/internal/log => ../../internal/log

replace github.com/gmax79/otusgolang/grpc/internal/calendar => ../../internal/calendar

replace github.com/gmax79/otusgolang/grpc/internal/support => ../../internal/support

replace github.com/gmax79/otusgolang/grpc/cmd/mycalendar/pbcalendar => ./pbcalendar

go 1.13

require (
	github.com/gmax79/otusgolang/grpc/cmd/mycalendar/pbcalendar v0.0.0-00010101000000-000000000000
	github.com/gmax79/otusgolang/grpc/internal/calendar v0.0.0-00010101000000-000000000000
	github.com/gmax79/otusgolang/grpc/internal/log v0.0.0-00010101000000-000000000000
	github.com/gmax79/otusgolang/grpc/internal/support v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.3.2
	go.uber.org/zap v1.13.0
	google.golang.org/grpc v1.25.1
)
