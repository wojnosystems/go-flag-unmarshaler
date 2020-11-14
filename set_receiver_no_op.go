package flag_unmarshaler

import into_struct "github.com/wojnosystems/go-into-struct"

type SetReceiverNoOp struct {
}

func (s *SetReceiverNoOp) ReceiveSet(fullPath into_struct.Path, flagName string, value string) {
	// do nothing
}
