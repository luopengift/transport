package transport

import (
	"errors"
)

var (
	// ErrClosedError ErrClosedError
	ErrClosedError = errors.New("Closed")
	// ErrOperationTimeoutError ErrOperationTimeoutError
	ErrOperationTimeoutError = errors.New("OperationTimeoutError: timeout")
	// ErrBufferClosedError ErrBufferClosedError
	ErrBufferClosedError = errors.New("BufferClosedError:chan is closed")
	// ErrReadBufferClosedError ErrReadBufferClosedError
	ErrReadBufferClosedError = errors.New("ReadBufferClosedError:chan is closed")
	// ErrWriteBufferClosedError ErrWriteBufferClosedError
	ErrWriteBufferClosedError = errors.New("WriteBufferClosedError:chan is closed")
	// ErrMaxBytesError ErrMaxBytesError
	ErrMaxBytesError = errors.New("MaxBytesError:message is larger than byte buffer")
	// ErrInputNullError ErrInputNullError
	ErrInputNullError = errors.New("InputError:plugin is null")
	// ErrOutputNullError ErrOutputNullError
	ErrOutputNullError = errors.New("OutputError:plugin is null")
	// ErrUnknownError ErrUnknownError
	ErrUnknownError = errors.New("UnknownError:unknow error")
)
