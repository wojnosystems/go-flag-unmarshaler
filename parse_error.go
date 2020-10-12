package go_flag_unmarshaller

import "fmt"

type ParseError struct {
	Path        StructFlagPath
	originalErr error
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("environment variable '%s' failed to parse '%s'", p.Path.FlagPath, p.originalErr.Error())
}
