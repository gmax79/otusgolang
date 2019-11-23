module github.com/gmax79/otusgolang/http/cmd/mycalendar

replace github.com/gmax79/otusgolang/http/internal/log => ../../internal/log

replace github.com/gmax79/otusgolang/http/internal/calendar => ../../internal/calendar

replace github.com/gmax79/otusgolang/http/internal/support => ../../internal/support

go 1.13

require (
	github.com/gmax79/otusgolang/http/internal/calendar v0.0.0-00010101000000-000000000000
	github.com/gmax79/otusgolang/http/internal/log v0.0.0-00010101000000-000000000000
	github.com/gmax79/otusgolang/http/internal/support v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.13.0
)
