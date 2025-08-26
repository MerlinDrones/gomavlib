package conversion

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDialect = `<?xml version="1.0"?>
<mavlink>
  <version>0</version>
  <dialect>0</dialect>
  <enums>
    <enum name="A_TYPE">
      <description>Detected Anomaly Types.</description>
      <entry value="0" name="A">
        <description>A.</description>
      </entry>
      <entry value="1" name="B">
        <description>B.</description>
      </entry>
      <entry value="2" name="C">
        <description>C.</description>
      </entry>
      <entry value="3" name="D">
        <description>D.</description>
      </entry>
      <entry value="4" name="E">
        <description>E</description>
      </entry>
      <entry value="1" name="BIT0" />
      <entry value="2**4" name="BIT4" />
      <entry value="0b000100000000" name="BIT8" />
      <entry value="0x10000" name="BIT16" />
      <entry value="0b1000000000000000000000000000000000000000000000000000000000000" name="BIT60" />
      <entry value="2305843009213693952" name="BIT61" />
      <entry value="2**62" name="BIT62" />
      <entry value="0x8000000000000000" name="BIT63" />
    </enum>
  </enums>
  <messages>
    <message id="43000" name="A_MESSAGE">
      <description>Detected anomaly info measured by onboard sensors and actuators. </description>
      <field type="uint8_t" name="test_uint8" enum="A_TYPE">a test uint8</field>
	  <field type="char[16]" name="Test_string">a test string</field>
	  <field type="uint32_t[4]" name="test_array">a test array</field>
	  <extensions/>
      <field type="uint8_t" name="mission_type" enum="MAV_MISSION_TYPE">a test extension</field>
    </message>
  </messages>
</mavlink>
`

func TestConversion(t *testing.T) {
	// Work in a temp dir so we don’t collide with local files.
	dir, err := os.MkdirTemp("", "gomavlib")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	prev, _ := os.Getwd()
	defer os.Chdir(prev)
	require.NoError(t, os.Chdir(dir))

	// Write the test dialect XML.
	require.NoError(t, os.WriteFile("testdialect.xml", []byte(testDialect), 0o644))

	// Generate concrete structs (link=false) so we can validate fields/types.
	require.NoError(t, Convert("testdialect.xml", false))

	// ---- Verify generated message ----
	msgPath := filepath.Join("testdialect", "message_a_message.go")
	buf, err := os.ReadFile(msgPath)
	require.NoError(t, err)
	gotMsg := string(buf)

	// Package and type
	require.Contains(t, gotMsg, "package testdialect")
	require.Contains(t, gotMsg, "type MessageAMessage struct")

	// Field names & types (don’t assert exact formatting or tags)
	require.Contains(t, gotMsg, "TestUint8 A_TYPE")
	require.Contains(t, gotMsg, "TestString string")
	require.Contains(t, gotMsg, "TestArray [4]uint32")
	require.Contains(t, gotMsg, "MissionType MAV_MISSION_TYPE")

	// GetID signature and value
	require.Contains(t, gotMsg, "func (*MessageAMessage) GetID() uint32")
	require.Contains(t, gotMsg, "return 43000")

	// ---- Verify generated enum ----
	enumPath := filepath.Join("testdialect", "enum_a_type.go")
	buf, err = os.ReadFile(enumPath)
	require.NoError(t, err)
	gotEnum := string(buf)

	// Package and base type (uint32 or uint64 depending on max value; BIT63 forces uint64)
	require.Contains(t, gotEnum, "package testdialect")
	require.Regexp(t, regexp.MustCompile(`type\s+A_TYPE\s+uint(32|64)`), gotEnum)

	// Constants (names and decimal values)
	require.Contains(t, gotEnum, `A A_TYPE = 0`)
	require.Contains(t, gotEnum, `B A_TYPE = 1`)
	require.Contains(t, gotEnum, `C A_TYPE = 2`)
	require.Contains(t, gotEnum, `D A_TYPE = 3`)
	require.Contains(t, gotEnum, `E A_TYPE = 4`)
	require.Contains(t, gotEnum, `BIT0 A_TYPE = 1`)
	require.Contains(t, gotEnum, `BIT4 A_TYPE = 16`)
	require.Contains(t, gotEnum, `BIT8 A_TYPE = 256`)
	require.Contains(t, gotEnum, `BIT16 A_TYPE = 65536`)
	require.Contains(t, gotEnum, `BIT60 A_TYPE = 1152921504606846976`)
	require.Contains(t, gotEnum, `BIT61 A_TYPE = 2305843009213693952`)
	require.Contains(t, gotEnum, `BIT62 A_TYPE = 4611686018427387904`)
	require.Contains(t, gotEnum, `BIT63 A_TYPE = 9223372036854775808`)

	// Methods exist (don’t care about internal map variable names)
	require.Contains(t, gotEnum, "func (e A_TYPE) MarshalText() ([]byte, error)")
	require.Contains(t, gotEnum, "func (e *A_TYPE) UnmarshalText(text []byte) error")
	require.Contains(t, gotEnum, "func (e A_TYPE) String() string")
}
