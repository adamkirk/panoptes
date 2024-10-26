package dt

import "time"

func NowUTC() time.Time {
	return time.Now().UTC()
}