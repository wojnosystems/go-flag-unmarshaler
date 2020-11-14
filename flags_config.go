package flag_unmarshaler

import parse_register "github.com/wojnosystems/go-parse-register"

type flagsConfig struct {
	// flags is the source of environment variables.
	// If you leave it blank, it will default to using the operating system environment variables with no prefixes.
	flags Reader
	// parseRegistry maps go-default and custom types to members of the provided structure. If left blank, defaults to just Go's primitives being mapped
	parseRegistry parse_register.ValueSetter
	emitter       SetReceiver
}

func newFlagsConfig(flags Reader, parseRegistry parse_register.ValueSetter, emitter SetReceiver) flagsConfig {
	if flags == nil {
		osGroup := SplitArgs()[0]
		flags = &osGroup
	}
	if parseRegistry == nil {
		parseRegistry = defaultParseRegister
	}
	if emitter == nil {
		emitter = defaultNoOpSetReceiver
	}
	return flagsConfig{
		flags:         flags,
		parseRegistry: parseRegistry,
		emitter:       emitter,
	}
}
