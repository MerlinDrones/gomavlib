//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// These defines are predefined OR-combined mode flags. There is no need to use values from this enum, but it
// simplifies the use of the mode flags. Note that manual input is enabled in all modes as a safety override.
type MAV_MODE uint64

const (
	// System is not ready to fly, booting, calibrating, etc. No flag is set.
	MAV_MODE_PREFLIGHT MAV_MODE = 0
	// System is allowed to be active, under assisted RC control.
	MAV_MODE_STABILIZE_DISARMED MAV_MODE = 80
	// System is allowed to be active, under assisted RC control.
	MAV_MODE_STABILIZE_ARMED MAV_MODE = 208
	// System is allowed to be active, under manual (RC) control, no stabilization
	MAV_MODE_MANUAL_DISARMED MAV_MODE = 64
	// System is allowed to be active, under manual (RC) control, no stabilization
	MAV_MODE_MANUAL_ARMED MAV_MODE = 192
	// System is allowed to be active, under autonomous control, manual setpoint
	MAV_MODE_GUIDED_DISARMED MAV_MODE = 88
	// System is allowed to be active, under autonomous control, manual setpoint
	MAV_MODE_GUIDED_ARMED MAV_MODE = 216
	// System is allowed to be active, under autonomous control and navigation (the trajectory is decided onboard and not pre-programmed by waypoints)
	MAV_MODE_AUTO_DISARMED MAV_MODE = 92
	// System is allowed to be active, under autonomous control and navigation (the trajectory is decided onboard and not pre-programmed by waypoints)
	MAV_MODE_AUTO_ARMED MAV_MODE = 220
	// UNDEFINED mode. This solely depends on the autopilot - use with caution, intended for developers only.
	MAV_MODE_TEST_DISARMED MAV_MODE = 66
	// UNDEFINED mode. This solely depends on the autopilot - use with caution, intended for developers only.
	MAV_MODE_TEST_ARMED MAV_MODE = 194
)

var labels_MAV_MODE = map[MAV_MODE]string{
	MAV_MODE_PREFLIGHT:          "MAV_MODE_PREFLIGHT",
	MAV_MODE_STABILIZE_DISARMED: "MAV_MODE_STABILIZE_DISARMED",
	MAV_MODE_STABILIZE_ARMED:    "MAV_MODE_STABILIZE_ARMED",
	MAV_MODE_MANUAL_DISARMED:    "MAV_MODE_MANUAL_DISARMED",
	MAV_MODE_MANUAL_ARMED:       "MAV_MODE_MANUAL_ARMED",
	MAV_MODE_GUIDED_DISARMED:    "MAV_MODE_GUIDED_DISARMED",
	MAV_MODE_GUIDED_ARMED:       "MAV_MODE_GUIDED_ARMED",
	MAV_MODE_AUTO_DISARMED:      "MAV_MODE_AUTO_DISARMED",
	MAV_MODE_AUTO_ARMED:         "MAV_MODE_AUTO_ARMED",
	MAV_MODE_TEST_DISARMED:      "MAV_MODE_TEST_DISARMED",
	MAV_MODE_TEST_ARMED:         "MAV_MODE_TEST_ARMED",
}

var values_MAV_MODE = map[string]MAV_MODE{
	"MAV_MODE_PREFLIGHT":          MAV_MODE_PREFLIGHT,
	"MAV_MODE_STABILIZE_DISARMED": MAV_MODE_STABILIZE_DISARMED,
	"MAV_MODE_STABILIZE_ARMED":    MAV_MODE_STABILIZE_ARMED,
	"MAV_MODE_MANUAL_DISARMED":    MAV_MODE_MANUAL_DISARMED,
	"MAV_MODE_MANUAL_ARMED":       MAV_MODE_MANUAL_ARMED,
	"MAV_MODE_GUIDED_DISARMED":    MAV_MODE_GUIDED_DISARMED,
	"MAV_MODE_GUIDED_ARMED":       MAV_MODE_GUIDED_ARMED,
	"MAV_MODE_AUTO_DISARMED":      MAV_MODE_AUTO_DISARMED,
	"MAV_MODE_AUTO_ARMED":         MAV_MODE_AUTO_ARMED,
	"MAV_MODE_TEST_DISARMED":      MAV_MODE_TEST_DISARMED,
	"MAV_MODE_TEST_ARMED":         MAV_MODE_TEST_ARMED,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_MODE) MarshalText() ([]byte, error) {
	if name, ok := labels_MAV_MODE[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_MODE) UnmarshalText(text []byte) error {
	if value, ok := values_MAV_MODE[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = MAV_MODE(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_MODE) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
