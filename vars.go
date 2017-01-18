package chain

import "errors"

var (
	ErrTypeNotFunction = errors.New("argument type not function")
	ErrInArgsMissMatch = errors.New("input arguments count not match")
	ErrOutCntMissMatch = errors.New("output parameter count not match")
	ErrExecuteTimeout  = errors.New("parallel execute timeout")
)
