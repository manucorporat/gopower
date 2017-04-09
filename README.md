# go-power
go-power allows developers to watch for power consumption of you system.

## Install

```
go get github.com/manucorporat/go-power
```

## Usage

### Instant values

```go
package main

import (
	"fmt"

	"github.com/manucorporat/gopower"
)

func main() {
	fmt.Println("Current:", gopower.CurrentNow())
	fmt.Println("Voltage:", gopower.VoltageNow())
	fmt.Println("Power:", gopower.PowerNow())
}
```

### Continuos monitoring

```go
package main

import (
	"fmt"

	"time"

	"github.com/manucorporat/gopower"
)

func main() {
	// Watcher that takes a sample each second and keeps a buffer for 1 minute.
	watcher := gopower.NewWatcher("power.log", time.Second, time.Minute)

	for range time.Tick(1 * time.Second) {
		sample := watcher.Mean(10 * time.Second)
		fmt.Println(sample, "\n")
	}
}

```

