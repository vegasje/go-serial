// +build !windows

package serial

// #include <termios.h>
// #include <unistd.h>
import "C"

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
)

func (c *Connection) open() (err error) {
	c.file, err = os.OpenFile(c.Name, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0666)
	if err != nil {
		return err
	}

	c.reader = bufio.NewReader(c.file)
	c.writer = bufio.NewWriter(c.file)

	fd := C.int(c.file.Fd())
	if C.isatty(fd) != 1 {
		c.file.Close()
		return errors.New("File is not a tty")
	}

	var st C.struct_termios
	_, err = C.tcgetattr(fd, &st)
	if err != nil {
		c.file.Close()
		return err
	}

	speed, err := convertBaud(c.Baud)
	if err != nil {
		c.file.Close()
		return err
	}

	_, err = C.cfsetispeed(&st, speed)
	if err != nil {
		c.file.Close()
		return err
	}

	_, err = C.cfsetospeed(&st, speed)
	if err != nil {
		c.file.Close()
		return err
	}

	// No timeout set directly on termios
	st.c_cc[C.VMIN] = 0
	st.c_cc[C.VTIME] = 0

	// Select local mode
	st.c_cflag |= (C.CLOCAL | C.CREAD)

	// Select raw mode
	st.c_lflag &= ^C.tcflag_t(C.ICANON | C.ECHO | C.ECHOE | C.ISIG)
	st.c_oflag &= ^C.tcflag_t(C.OPOST)

	_, err = C.tcsetattr(fd, C.TCSANOW, &st)
	if err != nil {
		c.file.Close()
		return err
	}

	r1, _, e := syscall.Syscall(syscall.SYS_FCNTL,
		uintptr(c.file.Fd()),
		uintptr(syscall.F_SETFL),
		uintptr(0))
	if e != 0 || r1 != 0 {
		s := fmt.Sprint("Clearing NONBLOCK syscall error:", e, r1)
		c.file.Close()
		return errors.New(s)
	}

	return nil
}

func (c *Connection) read(buf []byte) (size int, err error) {
	start := time.Now()

	if c.reader == nil {
		return size, errors.New("This connection has not been opened.")
	}

	available := c.reader.Buffered()

	for size < len(buf) {
		// Stop reading if we have reached the timeout
		current := time.Now()
		if current.Sub(start) >= c.Timeout {
			break
		}

		n, err := c.reader.Read(buf[size:available])
		if err != nil {
			return size, err
		}

		size += n
	}

	return size, nil
}

func convertBaud(baud Baud) (C.speed_t, error) {
	var speed C.speed_t

	switch baud {
	case BAUD_115200:
		return C.B115200, nil
	case BAUD_57600:
		return C.B57600, nil
	case BAUD_38400:
		return C.B38400, nil
	case BAUD_19200:
		return C.B19200, nil
	case BAUD_9600:
		return C.B9600, nil
	case BAUD_4800:
		return C.B4800, nil
	case BAUD_2400:
		return C.B2400, nil
	default:
		return speed, fmt.Errorf("Unknown baud rate: %v", baud)
	}
}
