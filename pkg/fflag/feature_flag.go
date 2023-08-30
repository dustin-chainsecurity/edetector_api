package fflag

import (
	"edetector_API/pkg/logger"

	flagsmith "github.com/Flagsmith/flagsmith-go-client"
)

var FFLAG *flagsmith.Client

func Get_fflag() {
	config := flagsmith.DefaultConfig()
	client := flagsmith.NewClient("ser.XSvF9amQyW6Pcamg8XxGHW", config)
	FFLAG = client
	isEnabled, _ := FFLAG.FeatureEnabled("always_true")
	if !isEnabled {
		logger.Error("Connection to Flagsmith failed")
	}
}