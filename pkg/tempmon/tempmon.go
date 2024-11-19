package tempmon

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v4/sensors"
)

// GetAverageCPUTemp retrieves the CPU temperature using preferred sensor logic.
func GetAverageCPUTemp() (float64, error) {
	stats, err := sensors.SensorsTemperatures()

	if err != nil {
		return 0, fmt.Errorf("failed to read sensor temperatures: %w", err)
	}

	var coreTemps []float64
	var preferredTemp float64

	for _, s := range stats {
		switch {
		// Use CPU package temps when available
		case s.SensorKey == "x86_pkg_temp" || s.SensorKey == "cpu_thermal":
			preferredTemp = s.Temperature
		// Or collect individual core temperatures otherwise
		case isCoreTempSensor(s.SensorKey):
			coreTemps = append(coreTemps, s.Temperature)
		}
	}

	if preferredTemp > 0 {
		return preferredTemp, nil
	}

	if preferredTemp == 0 && len(coreTemps) == 0 {
		return 0, fmt.Errorf("no CPU temperature sensors found")
	}

	var total float64

	for _, temp := range coreTemps {
		total += temp
	}

	average := total / float64(len(coreTemps))

	return average, nil
}

// isCoreTempSensor checks if a sensorKey likely represents a core temperature sensor
func isCoreTempSensor(sensorKey string) bool {
	sensorKey = strings.ToLower(sensorKey)
	corePatterns := []string{"coretemp", "cpu", "die"}

	for _, pattern := range corePatterns {
		if strings.Contains(sensorKey, pattern) {
			return true
		}
	}

	return false
}
