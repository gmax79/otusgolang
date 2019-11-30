module github.com/gmax79/otusgolang/grpc/cmd/mycalendar

replace github.com/gmax79/otusgolang/grpc/internal/log => ../../internal/log

replace github.com/gmax79/otusgolang/grpc/internal/calendar => ../../internal/calendar

replace github.com/gmax79/otusgolang/grpc/internal/support => ../../internal/support

go 1.13

require (
	github.com/gmax79/otusgolang/grpc/internal/calendar v0.0.0-00010101000000-000000000000
	github.com/gmax79/otusgolang/grpc/internal/log v0.0.0-00010101000000-000000000000
	github.com/gmax79/otusgolang/grpc/internal/support v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.13.0
)
