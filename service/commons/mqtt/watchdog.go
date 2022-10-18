package mqtt

import (
	"time"
)

type Watchdog struct {
	timer     *time.Timer
	timeout   time.Duration
	startChan chan bool
	stopChan  chan bool
}

// NewWatchdog ... ctor
func NewWatchdog(timeout time.Duration) *Watchdog {
	w := &Watchdog{
		timer:     time.NewTimer(timeout),
		timeout:   timeout,
		startChan: make(chan bool, 1),
		stopChan:  make(chan bool, 1),
	}
	w.timerExpireWorker()
	return w
}

func (w *Watchdog) timerExpireWorker() {
	go func(w *Watchdog) {
		for {
			select {
			case <-w.stopChan:
				return
			case <-w.timer.C:
				w.startChan <- true
				w.Reset()
			}
		}
	}(w)
}

func (w *Watchdog) Start() chan bool {
	return w.startChan
}

// Stop ... stops the watchdog
func (w *Watchdog) Stop() {
	w.timer.Stop()
	w.stopChan <- true
}

// Reset ... resets the timer
func (w *Watchdog) Reset() {
	w.timer.Stop()
	w.timer = time.NewTimer(w.timeout)
}
