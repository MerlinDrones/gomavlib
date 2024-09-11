//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Result from PARAM_EXT_SET message (or a PARAM_SET within a transaction).
type PARAM_ACK uint64

const (
	// Parameter value ACCEPTED and SET
	PARAM_ACK_ACCEPTED PARAM_ACK = 0
	// Parameter value UNKNOWN/UNSUPPORTED
	PARAM_ACK_VALUE_UNSUPPORTED PARAM_ACK = 1
	// Parameter failed to set
	PARAM_ACK_FAILED PARAM_ACK = 2
	// Parameter value received but not yet set/accepted. A subsequent PARAM_ACK_TRANSACTION or PARAM_EXT_ACK with the final result will follow once operation is completed. This is returned immediately for parameters that take longer to set, indicating that the the parameter was received and does not need to be resent.
	PARAM_ACK_IN_PROGRESS PARAM_ACK = 3
)

var labels_PARAM_ACK = map[PARAM_ACK]string{
	PARAM_ACK_ACCEPTED:          "PARAM_ACK_ACCEPTED",
	PARAM_ACK_VALUE_UNSUPPORTED: "PARAM_ACK_VALUE_UNSUPPORTED",
	PARAM_ACK_FAILED:            "PARAM_ACK_FAILED",
	PARAM_ACK_IN_PROGRESS:       "PARAM_ACK_IN_PROGRESS",
}

var values_PARAM_ACK = map[string]PARAM_ACK{
	"PARAM_ACK_ACCEPTED":          PARAM_ACK_ACCEPTED,
	"PARAM_ACK_VALUE_UNSUPPORTED": PARAM_ACK_VALUE_UNSUPPORTED,
	"PARAM_ACK_FAILED":            PARAM_ACK_FAILED,
	"PARAM_ACK_IN_PROGRESS":       PARAM_ACK_IN_PROGRESS,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e PARAM_ACK) MarshalText() ([]byte, error) {
	if name, ok := labels_PARAM_ACK[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *PARAM_ACK) UnmarshalText(text []byte) error {
	if value, ok := values_PARAM_ACK[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = PARAM_ACK(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e PARAM_ACK) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
