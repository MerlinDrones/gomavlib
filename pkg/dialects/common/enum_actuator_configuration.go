//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Actuator configuration, used to change a setting on an actuator. Component information metadata can be used to know which outputs support which commands.
type ACTUATOR_CONFIGURATION uint64

const (
	// Do nothing.
	ACTUATOR_CONFIGURATION_NONE ACTUATOR_CONFIGURATION = 0
	// Command the actuator to beep now.
	ACTUATOR_CONFIGURATION_BEEP ACTUATOR_CONFIGURATION = 1
	// Permanently set the actuator (ESC) to 3D mode (reversible thrust).
	ACTUATOR_CONFIGURATION_3D_MODE_ON ACTUATOR_CONFIGURATION = 2
	// Permanently set the actuator (ESC) to non 3D mode (non-reversible thrust).
	ACTUATOR_CONFIGURATION_3D_MODE_OFF ACTUATOR_CONFIGURATION = 3
	// Permanently set the actuator (ESC) to spin direction 1 (which can be clockwise or counter-clockwise).
	ACTUATOR_CONFIGURATION_SPIN_DIRECTION1 ACTUATOR_CONFIGURATION = 4
	// Permanently set the actuator (ESC) to spin direction 2 (opposite of direction 1).
	ACTUATOR_CONFIGURATION_SPIN_DIRECTION2 ACTUATOR_CONFIGURATION = 5
)

var labels_ACTUATOR_CONFIGURATION = map[ACTUATOR_CONFIGURATION]string{
	ACTUATOR_CONFIGURATION_NONE:            "ACTUATOR_CONFIGURATION_NONE",
	ACTUATOR_CONFIGURATION_BEEP:            "ACTUATOR_CONFIGURATION_BEEP",
	ACTUATOR_CONFIGURATION_3D_MODE_ON:      "ACTUATOR_CONFIGURATION_3D_MODE_ON",
	ACTUATOR_CONFIGURATION_3D_MODE_OFF:     "ACTUATOR_CONFIGURATION_3D_MODE_OFF",
	ACTUATOR_CONFIGURATION_SPIN_DIRECTION1: "ACTUATOR_CONFIGURATION_SPIN_DIRECTION1",
	ACTUATOR_CONFIGURATION_SPIN_DIRECTION2: "ACTUATOR_CONFIGURATION_SPIN_DIRECTION2",
}

var values_ACTUATOR_CONFIGURATION = map[string]ACTUATOR_CONFIGURATION{
	"ACTUATOR_CONFIGURATION_NONE":            ACTUATOR_CONFIGURATION_NONE,
	"ACTUATOR_CONFIGURATION_BEEP":            ACTUATOR_CONFIGURATION_BEEP,
	"ACTUATOR_CONFIGURATION_3D_MODE_ON":      ACTUATOR_CONFIGURATION_3D_MODE_ON,
	"ACTUATOR_CONFIGURATION_3D_MODE_OFF":     ACTUATOR_CONFIGURATION_3D_MODE_OFF,
	"ACTUATOR_CONFIGURATION_SPIN_DIRECTION1": ACTUATOR_CONFIGURATION_SPIN_DIRECTION1,
	"ACTUATOR_CONFIGURATION_SPIN_DIRECTION2": ACTUATOR_CONFIGURATION_SPIN_DIRECTION2,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e ACTUATOR_CONFIGURATION) MarshalText() ([]byte, error) {
	if name, ok := labels_ACTUATOR_CONFIGURATION[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *ACTUATOR_CONFIGURATION) UnmarshalText(text []byte) error {
	if value, ok := values_ACTUATOR_CONFIGURATION[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = ACTUATOR_CONFIGURATION(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e ACTUATOR_CONFIGURATION) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
