//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Optical flow from an angular rate flow sensor (e.g. PX4FLOW or mouse sensor)
type MessageOpticalFlowRad struct {
	// Timestamp (UNIX Epoch time or time since system boot). The receiving end can infer timestamp format (since 1.1.1970 or since system boot) by checking for the magnitude of the number.
	TimeUsec uint64
	// Sensor ID
	SensorId uint8
	// Integration time. Divide integrated_x and integrated_y by the integration time to obtain average flow. The integration time also indicates the.
	IntegrationTimeUs uint32
	// Flow around X axis (Sensor RH rotation about the X axis induces a positive flow. Sensor linear motion along the positive Y axis induces a negative flow.)
	IntegratedX float32
	// Flow around Y axis (Sensor RH rotation about the Y axis induces a positive flow. Sensor linear motion along the positive X axis induces a positive flow.)
	IntegratedY float32
	// RH rotation around X axis
	IntegratedXgyro float32
	// RH rotation around Y axis
	IntegratedYgyro float32
	// RH rotation around Z axis
	IntegratedZgyro float32
	// Temperature
	Temperature int16
	// Optical flow quality / confidence. 0: no valid flow, 255: maximum quality
	Quality uint8
	// Time since the distance was sampled.
	TimeDeltaDistanceUs uint32
	// Distance to the center of the flow field. Positive value (including zero): distance known. Negative value: Unknown distance.
	Distance float32
}

// GetID implements the message.Message interface.
func (*MessageOpticalFlowRad) GetID() uint32 {
	return 106
}
