package flag_unmarshaler

import (
	into_struct "github.com/wojnosystems/go-into-struct"
	"github.com/wojnosystems/go-parse-register"
)

type Flags struct {
	config flagsConfig
}

func New(flags Reader) *Flags {
	return NewWithTypeParsers(
		flags,
		defaultParseRegister)
}

func NewWithTypeParsers(flags Reader, parseRegistry parse_register.ValueSetter) *Flags {
	return &Flags{
		config: newFlagsConfig(
			flags,
			parseRegistry,
			defaultNoOpSetReceiver,
		),
	}
}

func NewWithEmitter(flags Reader, emitter SetReceiver) *Flags {
	return &Flags{
		config: newFlagsConfig(
			flags,
			defaultParseRegister,
			emitter,
		),
	}
}

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func (e *Flags) Unmarshal(into interface{}) (err error) {
	return into_struct.Unmarshall(into, newFlagsInternal(e.config, into))
}
