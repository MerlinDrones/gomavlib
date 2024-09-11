//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Data for injecting into the onboard GPS (used for DGPS)
type MessageGpsInjectData struct {
	// System ID
	TargetSystem uint8
	// Component ID
	TargetComponent uint8
	// Data length
	Len uint8
	// Raw data (110 is enough for 12 satellites of RTCMv2)
	Data [110]uint8
}

// GetID implements the message.Message interface.
func (*MessageGpsInjectData) GetID() uint32 {
	return 123
}
