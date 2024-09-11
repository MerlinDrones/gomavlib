//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Illuminator status
type MessageIlluminatorStatus struct {
	// Time since the start-up of the illuminator in ms
	UptimeMs uint32
	// 0: Illuminators OFF, 1: Illuminators ON
	Enable uint8
	// Supported illuminator modes
	ModeBitmask ILLUMINATOR_MODE `mavenum:"uint8"`
	// Errors
	ErrorStatus ILLUMINATOR_ERROR_FLAGS `mavenum:"uint32"`
	// Illuminator mode
	Mode ILLUMINATOR_MODE `mavenum:"uint8"`
	// Illuminator brightness
	Brightness float32
	// Illuminator strobing period in seconds
	StrobePeriod float32
	// Illuminator strobing duty cycle
	StrobeDutyCycle float32
	// Temperature in Celsius
	TempC float32
	// Minimum strobing period in seconds
	MinStrobePeriod float32
	// Maximum strobing period in seconds
	MaxStrobePeriod float32
}

// GetID implements the message.Message interface.
func (*MessageIlluminatorStatus) GetID() uint32 {
	return 440
}
