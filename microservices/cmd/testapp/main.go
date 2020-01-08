package main

import (
	"fmt"
	"net/http"

	tests "github.com/gmax79/otusgolang/microservices/internal/stests"
)

const host = "http://localhost:8888"

func main() {
	fmt.Println("Testing calendar app")
	tests.PostWithPrint(host, "", map[string]string{}, http.StatusNotFound)
	tests.PostWithPrint(host, "a", map[string]string{}, http.StatusNotFound)
	tests.PostWithPrint(host, "b", map[string]string{}, http.StatusNotFound)

	r1 := map[string]string{
		"time":  "2020-10-22 18:00:00",
		"event": "Maks birthday",
	}
	tests.PostWithPrint(host, "create_event", r1, http.StatusOK)

	r2 := map[string]string{
		"time":  "2020-10-22 18:00:00",
		"event": "Maks birthday",
	}
	tests.PostWithPrint(host, "delete_event", r2, http.StatusOK)

	r3 := map[string]string{
		"time":  "2020-03-07 12:00:00",
		"event": "Party",
	}
	tests.PostWithPrint(host, "create_event", r3, http.StatusOK)

	r3old := map[string]string{
		"time":  "2020-03-10 18:00:00",
		"event": "Party",
	}
	tests.PostWithPrint(host, "delete_event", r3old, http.StatusOK)

	r3move := map[string]string{
		"time":    "2020-03-07 12:00:00",
		"event":   "Party",
		"newtime": "2020-03-10 18:00:00",
	}
	tests.PostWithPrint(host, "move_event", r3move, http.StatusOK)

	r5 := map[string]string{
		"time":  "2020-03-08 16:00:00",
		"event": "Party in club",
	}
	tests.PostWithPrint(host, "create_event", r5, http.StatusOK)

	r6 := map[string]string{
		"time":  "2020-03-15 12:00:00",
		"event": "Exam",
	}
	tests.PostWithPrint(host, "create_event", r6, http.StatusOK)

	tests.GetWithPrint(host, "events_for_day?day=2020-03-07", http.StatusOK, 0)
	tests.GetWithPrint(host, "events_for_day?day=2020-03-10", http.StatusOK, 1)
	tests.GetWithPrint(host, "events_for_week?week=2020-11", http.StatusOK, 2)
	tests.GetWithPrint(host, "events_for_month?month=2020-03", http.StatusOK, 3)

	fmt.Println("Tests via http interface OK!")
}
