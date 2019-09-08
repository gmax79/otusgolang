package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

const ntpServer = "0.beevik-ntp.pool.ntp.org"

func main() {
	timeval, err := ntp.Time(ntpServer)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	formatedTime := timeval.Format(time.RFC822)
	fmt.Println(formatedTime)
}
