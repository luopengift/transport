package tcp

import (
	"crypto/tls"
	"fmt"
	"github.com/luopengift/transport"
	"io"
	"log"
	"net"
	"time"
)

// TcpOutput used for sending raw tcp payloads
// Currently used for internal communication between listener and replay server
// Can be used for transfering binary payloads like protocol buffers
type TcpOutput struct {
	address string
	limit   int
	buf     chan []byte
	config  *TcpOutputConfig
}

type TcpOutputConfig struct {
	secure bool
}

// NewTcpOutput constructor for TcpOutput
// Initialize 10 workers which hold keep-alive connection
func NewTcpOutput() *TcpOutput {
	o := new(TcpOutput)
	o.buf = make(chan []byte, 100)

	for i := 0; i < 10; i++ {
		go o.worker()
	}

	return o
}

func (o *TcpOutput) worker() {
	retries := 1
	conn, err := o.connect(o.address)
	for {
		if err == nil {
			break
		}

		log.Println("Can't connect to aggregator instance, reconnecting in 1 second. Retries:", retries)
		time.Sleep(1 * time.Second)

		conn, err = o.connect(o.address)
		retries++
	}

	if retries > 0 {
		log.Println("Connected to aggregator instance after ", retries, " retries")
	}

	defer conn.Close()

	for {
		conn.Write(<-o.buf)
		_, err := conn.Write([]byte(payloadSeparator))

		if err != nil {
			log.Println("Lost connection with aggregator instance, reconnecting")
			go o.worker()
			break
		}
	}
}

func (o *TcpOutput) Write(data []byte) (n int, err error) {
	if !isOriginPayload(data) {
		return len(data), nil
	}

	// We have to copy, because sending data in multiple threads
	newBuf := make([]byte, len(data))
	copy(newBuf, data)

	o.buf <- newBuf

	return len(data), nil
}

func (o *TcpOutput) connect(address string) (conn net.Conn, err error) {
	if o.config.secure {
		conn, err = tls.Dial("tcp", address, &tls.Config{})
	} else {
		conn, err = net.Dial("tcp", address)
	}

	return
}

func (o *TcpOutput) Close() error {
	return nil
}

func (o *TcpOutput) String() string {
	return fmt.Sprintf("TCP output %s, limit: %d", o.address, o.limit)
}

func init() {
	transport.RegistOutputer("tcp", NewTcpOutput())
}
