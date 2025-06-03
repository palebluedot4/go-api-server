package timeutil

import (
	"time"
)

var TaipeiLocation *time.Location

func init() {
	var err error
	TaipeiLocation, err = time.LoadLocation("Asia/Taipei")
	if err != nil {
		TaipeiLocation = time.FixedZone("Asia/Taipei", 8*60*60)
	}
}

func Now() time.Time {
	return time.Now().In(TaipeiLocation)
}
