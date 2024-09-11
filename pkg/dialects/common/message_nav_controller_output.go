//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// The state of the navigation and position controller.
type MessageNavControllerOutput struct {
	// Current desired roll
	NavRoll float32
	// Current desired pitch
	NavPitch float32
	// Current desired heading
	NavBearing int16
	// Bearing to current waypoint/target
	TargetBearing int16
	// Distance to active waypoint
	WpDist uint16
	// Current altitude error
	AltError float32
	// Current airspeed error
	AspdError float32
	// Current crosstrack error on x-y plane
	XtrackError float32
}

// GetID implements the message.Message interface.
func (*MessageNavControllerOutput) GetID() uint32 {
	return 62
}
