package utils

import "time"

type Time interface {
	Now() time.Time
}

type TimeImpl struct{}

func NewTime() Time {
	return &TimeImpl{}
}

func (t *TimeImpl) Now() time.Time {
	return time.Now()
}
