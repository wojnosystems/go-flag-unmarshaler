package go_flag_unmarshaller

type SetReceiverNoOp struct {
}

func (s *SetReceiverNoOp) ReceiveSet(fullPath string, envName string, value string) {
	// do nothing
}
