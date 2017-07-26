package pipeline

import (
	"errors"
)

var (
    ClosedError = errors.New("Closed")

    ReadBufferClosedError = errors.New("Read Buffer is closed")
    WriteBufferClosedError = errors.New("Write Buffer is closed")

	MaxBytesError   = errors.New("MaxBytesError:message is larger than byte buffer")
	InputNullError  = errors.New("InputError:plugin is null")
	OutputNullError = errors.New("OutputError:plugin is null")
	UnknownError    = errors.New("UnknownError:unknow error")
)
