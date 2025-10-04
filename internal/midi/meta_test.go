package midi

import (
	"encoding/binary"
	"testing"
)

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic but function did not panic")
		}
	}()
	f()
}

func TestThreeByteIntConversion(t *testing.T) {
	t.Helper()
	data := []byte{0x01, 0x02, 0x03} // Should be 0x010203 == 66051

	mustPanic(t, func() {
		// This will not give the correct value
		val := binary.BigEndian.Uint32(data) // Only first byte
		if val == 66051 {
			t.Errorf("Incorrect cast should not produce correct value")
		}
	})

	t.Run("correct conversion", func(t *testing.T) {
		val := int(data[0])<<16 | int(data[1])<<8 | int(data[2])
		if val != 66051 {
			t.Errorf("Expected 66051, got %d", val)
		}
	})
}
