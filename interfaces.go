package flag_unmarshaller

// flagReader reads environment variables
type flagReader interface {
	// Get the value of a single flag with any of the names in name
	Get(flagNamed string) (value string, ok bool)
	// Keys get a list of keys that begin with the prefix. If "" is passed, matches all and returns all keys
	Keys(prefix string) []string
}

type SetReceiver interface {
	// Receive the notice that a value was parsed and set at the fullPath in the destination structure
	// This will allow the flick library to know which values were updated from which source.
	ReceiveSet(structPath string, flagName string, value string)
}

type Unmarshaller interface {
	Unmarshall(into interface{}) error
	UnmarshallWithEmitter(into interface{}, emitter SetReceiver) error
}
