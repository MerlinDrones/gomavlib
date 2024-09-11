//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Tune formats supported by vehicle. This should be emitted as response to MAV_CMD_REQUEST_MESSAGE.
type MessageSupportedTunes struct {
	// System ID
	TargetSystem uint8
	// Component ID
	TargetComponent uint8
	// Bitfield of supported tune formats.
	Format TUNE_FORMAT `mavenum:"uint32"`
}

// GetID implements the message.Message interface.
func (*MessageSupportedTunes) GetID() uint32 {
	return 401
}
