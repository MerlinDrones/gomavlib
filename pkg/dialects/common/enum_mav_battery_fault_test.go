//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAV_BATTERY_FAULT(t *testing.T) {
	var e MAV_BATTERY_FAULT
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
