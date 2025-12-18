package types

import "time"

func IsLocked(currentTime time.Time, unlockDate string) bool {
	unlockTime, err := time.Parse(time.DateOnly, unlockDate)
	if err != nil {
		return false
	}

	currentDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
	return currentDay.Before(unlockTime)
}
