//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Sent from simulation to autopilot. The RAW values of the RC channels received. The standard PPM modulation is as follows: 1000 microseconds: 0%, 2000 microseconds: 100%. Individual receivers/transmitters might violate this specification.
type MessageHilRcInputsRaw struct {
	// Timestamp (UNIX Epoch time or time since system boot). The receiving end can infer timestamp format (since 1.1.1970 or since system boot) by checking for the magnitude of the number.
	TimeUsec uint64
	// RC channel 1 value
	Chan1Raw uint16
	// RC channel 2 value
	Chan2Raw uint16
	// RC channel 3 value
	Chan3Raw uint16
	// RC channel 4 value
	Chan4Raw uint16
	// RC channel 5 value
	Chan5Raw uint16
	// RC channel 6 value
	Chan6Raw uint16
	// RC channel 7 value
	Chan7Raw uint16
	// RC channel 8 value
	Chan8Raw uint16
	// RC channel 9 value
	Chan9Raw uint16
	// RC channel 10 value
	Chan10Raw uint16
	// RC channel 11 value
	Chan11Raw uint16
	// RC channel 12 value
	Chan12Raw uint16
	// Receive signal strength indicator in device-dependent units/scale. Values: [0-254], UINT8_MAX: invalid/unknown.
	Rssi uint8
}

// GetID implements the message.Message interface.
func (*MessageHilRcInputsRaw) GetID() uint32 {
	return 92
}
