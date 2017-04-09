package gopower

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Watcher allows you, to monitor your electric consumption in an easy and
// efficient way.
type Watcher struct {
	log      io.Writer
	interval time.Duration
	samples  []Sample
	index    int
	length   int
	mutex    sync.RWMutex
}

// Result contains the average of several samples given some requirements.
type Result struct {
	Sample

	TimeFrame time.Duration
	NuSamples int
}

func (r Result) String() string {
	return fmt.Sprintf(`# samples: %d
TimeFrame: %v
Current: %v
Voltage: %v
Power: %v`, r.NuSamples, r.TimeFrame, r.Current, r.Voltage, r.Power)
}

// NewWatcher creates a new Watcher given a sample interval and the buffer size
// specified in time.
func NewWatcher(logfile string, interval time.Duration, maxTime time.Duration) *Watcher {
	var log *os.File
	var err error
	if logfile == "" {
		log, err = os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		must(err)
	}

	bufferSize := maxTime / interval
	if bufferSize < 2 {
		panic("buffer is too small")
	}
	watcher := &Watcher{
		log:      log,
		interval: interval,
		index:    0,
		samples:  make([]Sample, bufferSize, bufferSize),
	}
	go watcher.start()
	return watcher
}

// Mean returns the average consumption of your system during the specified duration.
func (w *Watcher) Mean(duration time.Duration) Result {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	var sumA Ampere
	var sumV Volt
	var sumW Watt
	var lastSample time.Time

	nu := 0
	i := w.index - 1
	if i < 0 {
		return Result{}
	}
	start := w.samples[i].Instant
	for {
		sample := w.samples[i]
		if start.Sub(sample.Instant) > duration {
			break
		}
		if nu >= cap(w.samples) {
			break
		}
		lastSample = sample.Instant
		sumA += sample.Current
		sumV += sample.Voltage
		sumW += sample.Power
		nu++

		if i == 0 {
			i = cap(w.samples)
		}
		i--
	}

	return Result{
		Sample: Sample{
			Instant: time.Now(),
			Current: sumA / Ampere(nu),
			Voltage: sumV / Volt(nu),
			Power:   sumW / Watt(nu),
		},
		TimeFrame: start.Sub(lastSample),
		NuSamples: nu,
	}
}

func (w *Watcher) start() {
	for range time.Tick(w.interval) {
		sample, err := GetNow()
		if err != nil {
			log.Println("ERROR gopower:", err)
			continue
		}
		w.addSample(sample)
	}
}

func (w *Watcher) addSample(sample Sample) {
	w.mutex.Lock()
	w.samples[w.index] = sample

	if w.log != nil {
		fmt.Fprintf(w.log, "%d,%.0f,%.0f,%.0f\n",
			sample.Instant.UnixNano(),
			sample.Current,
			sample.Voltage,
			sample.Power)
	}

	// Circular buffer
	w.index++
	if w.index > cap(w.samples) {
		w.index = 0
	}
	w.mutex.Unlock()
}
