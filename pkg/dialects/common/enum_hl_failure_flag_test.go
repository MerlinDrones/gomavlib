//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_HL_FAILURE_FLAG(t *testing.T) {
	var e HL_FAILURE_FLAG
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
