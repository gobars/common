package dates

import (
	"time"
)

func ToDateStr() string {
	currentTime := time.Now()
	return currentTime.Format("20060102")
}
