package templates

import (
	"text/template"
)

var TplTest = template.Must(template.New("").Parse(
	`//autogenerated:yes
//nolint:revive
package dialects

import (
	"testing"
	"encoding"
	"reflect"

	"github.com/stretchr/testify/require"

	"github.com/merlindrones/gomavlib/pkg/dialects/common"
)

var casesEnum = []struct {
	name string
	dec  encoding.TextMarshaler
	enc  string
}{
	{
		"bitmask",
		common.POSITION_TARGET_TYPEMASK_VX_IGNORE | common.POSITION_TARGET_TYPEMASK_VY_IGNORE,
		"POSITION_TARGET_TYPEMASK_VX_IGNORE | POSITION_TARGET_TYPEMASK_VY_IGNORE",
	},
	{
		"value",
		common.GPS_FIX_TYPE_NO_FIX,
		"GPS_FIX_TYPE_NO_FIX",
	},
}

func TestEnumUnmarshalText(t *testing.T) {
	for _, ca := range casesEnum {
		t.Run(ca.name, func(t *testing.T) {
			dec := reflect.New(reflect.TypeOf(ca.dec)).Interface().(encoding.TextUnmarshaler)
			err := dec.UnmarshalText([]byte(ca.enc))
			require.NoError(t, err)
			require.Equal(t, ca.dec, reflect.ValueOf(dec).Elem().Interface())
		})
	}
}

func TestEnumMarshalText(t *testing.T) {
	for _, ca := range casesEnum {
		t.Run(ca.name, func(t *testing.T) {
			byts, err := ca.dec.MarshalText()
			require.NoError(t, err)
			require.Equal(t, ca.enc, string(byts))
		})
	}
}
`))

var TplDialectTest = template.Must(template.New("").Parse(
	`//autogenerated:yes
//nolint:revive
package {{ .PkgName }}

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/merlindrones/gomavlib/pkg/dialect"
)

func TestDialect(t *testing.T) {
	_, err := dialect.NewReadWriter(Dialect)
	require.NoError(t, err)
}
`))

var TplEnumTest = template.Must(template.New("").Parse(
	`//autogenerated:yes
//nolint:revive,govet,errcheck
package {{ .PkgName }}

import (
	"testing"
)

func TestEnum_{{ .Name }}(t *testing.T) {
	var e {{ .Name }}
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
`))
