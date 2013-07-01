package serial

import (
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	c, err := Open("/dev/ttyAMA0", 0x0, 500*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Write([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 128)
	_, err = s.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
}
