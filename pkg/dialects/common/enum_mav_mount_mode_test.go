//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAV_MOUNT_MODE(t *testing.T) {
	var e MAV_MOUNT_MODE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
