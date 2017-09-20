package transport

import (
	"errors"
)

var (
	ClosedError = errors.New("Closed")

	BufferClosedError      = errors.New("BufferClosedError:chan is closed")
	ReadBufferClosedError  = errors.New("ReadBufferClosedError:chan is closed")
	WriteBufferClosedError = errors.New("WriteBufferClosedError:chan is closed")

	MaxBytesError   = errors.New("MaxBytesError:message is larger than byte buffer")
	InputNullError  = errors.New("InputError:plugin is null")
	OutputNullError = errors.New("OutputError:plugin is null")
	UnknownError    = errors.New("UnknownError:unknow error")
)
