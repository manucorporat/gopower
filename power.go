package gopower

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	units "github.com/docker/go-units"
)

var ampUnits = []string{"uA", "mA", "A"}
var voltUnits = []string{"uV", "mV", "V"}
var wattUnits = []string{"uW", "mW", "W"}

// Ampere represents the Micro Amperes. SI Unit for electric current.
type Ampere float64

// Volt represents micro volts. SI Unit for electric potencial.
type Volt float64

// Watt represents micro watts. SI Unit for power.
type Watt float64

func (a Ampere) String() string {
	return units.CustomSize("%.3f%s", float64(a), 1000, ampUnits)
}

func (v Volt) String() string {
	return units.CustomSize("%.3f%s", float64(v), 1000, voltUnits)
}

func (w Watt) String() string {
	return units.CustomSize("%.3f%s", float64(w), 1000, wattUnits)
}

// Sample stores the current, voltage and power consumption of your system
// in a specific instant.
type Sample struct {
	Instant time.Time
	Current Ampere
	Voltage Volt
	Power   Watt
}

func readNumber(filepath string) (float64, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return 0, err
	}
	strData := strings.TrimSpace(string(data))
	value, err := strconv.ParseFloat(strData, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// GetCurrentNow returns the instant current of your system.
// Otherwise it returns an error.
func GetCurrentNow() (Ampere, error) {
	value, err := readNumber("/sys/class/power_supply/BAT0/current_now")
	if err != nil {
		return 0, err
	}
	return Ampere(value), nil
}

// GetVoltageNow returns the instant voltage of your system.
// Otherwise it returns an error.
func GetVoltageNow() (Volt, error) {
	value, err := readNumber("/sys/class/power_supply/BAT0/voltage_now")
	if err != nil {
		return 0, err
	}
	return Volt(value), nil
}

// GetNow returns a Sample. Ie. Instant current, voltage and power of your system.
// Otherwise it returns an error.
func GetNow() (Sample, error) {
	amperes, err := GetCurrentNow()
	if err != nil {
		return Sample{}, err
	}
	volts, err := GetVoltageNow()
	if err != nil {
		return Sample{}, err
	}
	watts := Watt(float64(amperes) * float64(volts) / 1e6)

	return Sample{
		Instant: time.Now(),
		Current: amperes,
		Voltage: volts,
		Power:   watts,
	}, nil
}

// GetPowerNow returns the instant power of your system.
// Otherwise it returns an error.
func GetPowerNow() (Watt, error) {
	sample, err := GetNow()
	if err != nil {
		return 0, err
	}
	return sample.Power, nil
}

// CurrentNow returns the instant current of your system.
// Otherwise it panics
func CurrentNow() Ampere {
	amperage, err := GetCurrentNow()
	must(err)
	return amperage
}

// VoltageNow returns the instant voltage of your system.
// Otherwise it panics
func VoltageNow() Volt {
	volts, err := GetVoltageNow()
	must(err)
	return volts
}

// PowerNow returns the instant power of your system.
// Otherwise it panics
func PowerNow() Watt {
	watts, err := GetPowerNow()
	must(err)
	return watts
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
