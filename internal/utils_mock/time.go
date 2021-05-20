package utils_mock

import (
	time "time"
)

type TimeTravelingMock struct {
	CurrentTime time.Time
}

func NewTimeTravelingMock(startTime time.Time) *TimeTravelingMock {
	return &TimeTravelingMock{CurrentTime: startTime}
}

func (t *TimeTravelingMock) Now() time.Time {
	return t.CurrentTime
}

func (t *TimeTravelingMock) ForwardBy(duration time.Duration) {
	t.CurrentTime.Add(duration)
}
