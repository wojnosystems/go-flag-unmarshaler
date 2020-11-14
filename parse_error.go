package flag_unmarshaler

import "fmt"

type ParseError struct {
	Path        StructFlagPath
	originalErr error
}

func newParseError(structFullPath, flagPath string, original error) *ParseError {
	return &ParseError{
		Path: StructFlagPath{
			StructPath: structFullPath,
			FlagPath:   flagPath,
		},
		originalErr: original,
	}
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("flag '%s' failed to parse because %s", p.Path.FlagPath, p.originalErr.Error())
}
