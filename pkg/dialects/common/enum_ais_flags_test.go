//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_AIS_FLAGS(t *testing.T) {
	var e AIS_FLAGS
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
