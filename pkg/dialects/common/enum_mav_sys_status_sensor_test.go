//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAV_SYS_STATUS_SENSOR(t *testing.T) {
	var e MAV_SYS_STATUS_SENSOR
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
