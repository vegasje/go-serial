package serial

import (
	"errors"
	"io"
	"os"
	"time"
)

const (
	BAUD_115200 Baud = 115200
	BAUD_57600  Baud = 57600
	BAUD_38400  Baud = 38400
	BAUD_19200  Baud = 19200
	BAUD_9600   Baud = 9600
	BAUD_4800   Baud = 4800
	BAUD_2400   Baud = 2400
)

type Baud int

type Connection struct {
	Name    string
	Baud    Baud
	Timeout time.Duration
	file    *os.File
}

func Open(name string, baud Baud, timeout time.Duration) (*Connection, error) {
	c := Connection{Name: name, Baud: baud, Timeout: timeout}

	err := c.open()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Connection) SetBaudRate(baud Baud) error {
	return c.setBaudRate(baud)
}

// Read size bytes from the serial port.
// If a timeout has been set > 0 it may return less characters than requested.
// With a timeout = 0 it will block until the requested number of bytes is read.
func (c *Connection) Read(buf []byte) (size int, err error) {
	start := time.Now()

	if c.file == nil {
		return size, errors.New("This connection has not been opened.")
	}

	for size < len(buf) {
		if c.Timeout != 0 {
			// Stop reading if we have reached the timeout
			current := time.Now()
			if current.Sub(start) >= c.Timeout {
				break
			}
		}

		n, err := c.file.Read(buf[size:])
		if err != nil && err != io.EOF {
			return size, err
		}

		size += n
	}

	return size, nil
}

func (c *Connection) Write(data []byte) (int, error) {
	return c.file.Write(data)
}

func (c *Connection) Close() error {
	return c.file.Close()
}
