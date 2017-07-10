package transport


import (
	"errors"
)

var (
	MaxBytesError = errors.New("MaxBytesError:message is larger than byte buffer")
	InputError = errors.New("InputError:plugin is null")
	OutputError = errors.New("OutputError:plugin is null")
	UnknownError = errors.New("UnknownError:unknow error")
)


