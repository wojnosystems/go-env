package env_parser

import into_struct "github.com/wojnosystems/go-into-struct"

type SetReceiverNoOp struct {
}

func (s *SetReceiverNoOp) ReceiveSet(fullPath into_struct.Path, envName string, value string) {
	// do nothing
}
