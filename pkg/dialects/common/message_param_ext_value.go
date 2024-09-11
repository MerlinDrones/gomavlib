//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Emit the value of a parameter. The inclusion of param_count and param_index in the message allows the recipient to keep track of received parameters and allows them to re-request missing parameters after a loss or timeout.
type MessageParamExtValue struct {
	// Parameter id, terminated by NULL if the length is less than 16 human-readable chars and WITHOUT null termination (NULL) byte if the length is exactly 16 chars - applications have to provide 16+1 bytes storage if the ID is stored as string
	ParamId string `mavlen:"16"`
	// Parameter value
	ParamValue string `mavlen:"128"`
	// Parameter type.
	ParamType MAV_PARAM_EXT_TYPE `mavenum:"uint8"`
	// Total number of parameters
	ParamCount uint16
	// Index of this parameter
	ParamIndex uint16
}

// GetID implements the message.Message interface.
func (*MessageParamExtValue) GetID() uint32 {
	return 322
}
