//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Camera absolute thermal range. This can be streamed when the associated VIDEO_STREAM_STATUS `flag` field bit VIDEO_STREAM_STATUS_FLAGS_THERMAL_RANGE_ENABLED is set, but a GCS may choose to only request it for the current active stream. Use MAV_CMD_SET_MESSAGE_INTERVAL to define message interval (param3 indicates the stream id of the current camera, or 0 for all streams, param4 indicates the target camera_device_id for autopilot-attached cameras or 0 for MAVLink cameras).
type MessageCameraThermalRange struct {
	// Timestamp (time since system boot).
	TimeBootMs uint32
	// Video Stream ID (1 for first, 2 for second, etc.)
	StreamId uint8
	// Camera id of a non-MAVLink camera attached to an autopilot (1-6).  0 if the component is a MAVLink camera (with its own component id).
	CameraDeviceId uint8
	// Temperature max.
	Max float32
	// Temperature max point x value (normalized 0..1, 0 is left, 1 is right), NAN if unknown.
	MaxPointX float32
	// Temperature max point y value (normalized 0..1, 0 is top, 1 is bottom), NAN if unknown.
	MaxPointY float32
	// Temperature min.
	Min float32
	// Temperature min point x value (normalized 0..1, 0 is left, 1 is right), NAN if unknown.
	MinPointX float32
	// Temperature min point y value (normalized 0..1, 0 is top, 1 is bottom), NAN if unknown.
	MinPointY float32
}

// GetID implements the message.Message interface.
func (*MessageCameraThermalRange) GetID() uint32 {
	return 277
}
