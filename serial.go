package serial

import (
	"bufio"
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
	reader  *bufio.Reader
	writer  *bufio.Writer
}

func Open(name string, baud Baud, timeout time.Duration) (*Connection, error) {
	c := Connection{Name: name, Baud: baud, Timeout: timeout}

	err := c.open()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Read size bytes from the serial port.
// If a timeout is set it may return less characters than requested.
// With no timeout it will block until the requested number of bytes is read.
func (c *Connection) Read(buf []byte) (int, error) {
	return c.read(buf)
}

func (c *Connection) Write(data []byte) (int, error) {
	return c.writer.Write(data)
}

func (c *Connection) Close() error {
	c.writer.Flush()
	return c.file.Close()
}
