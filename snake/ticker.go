package snake

import (
	"time"
)

type Ticker struct {
	interval    time.Duration
	min         time.Duration
	Coefficient float32
}

func NewTicker(initInterval time.Duration, min time.Duration, coefficient float32) *Ticker {
	return &Ticker{
		interval:    initInterval,
		min:         min,
		Coefficient: coefficient,
	}
}

func (t *Ticker) Interval() time.Duration {
	return t.interval
}

func (t *Ticker) ReduceInterval() {
	interval := time.Duration(float32(t.interval) * t.Coefficient)
	t.interval = min(interval, t.min)
}
