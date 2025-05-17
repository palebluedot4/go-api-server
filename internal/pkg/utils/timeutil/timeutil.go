package timeutil

import (
	"time"
)

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "2006-01-02 15:04:05"
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
