package agent

import (
	"flag"
)

type Settings struct {
	serverHost      string
	pollIntervalSec int
	reportInterval  int
}

var settings Settings

func Init() Settings {
	flag.StringVar(&settings.serverHost, "a", "localhost:8080", "server host")
	flag.IntVar(&settings.pollIntervalSec, "p", 2, "frequency of sending metrics to the server")
	flag.IntVar(&settings.reportInterval, "r", 10, "frequency of polling metrics from the runtime package")
	flag.Parse()
	if len(flag.Args()) > 0 {
		panic("used not declared arguments")
	}
	return settings
}
