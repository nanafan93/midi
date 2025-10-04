package vlq

import (
	"bytes"
	"io"
	"testing"
)

func TestParseMidiVLQ(t *testing.T) {
	type args struct {
		input io.ByteReader
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{
			"0x00 == 0x00",
			args{bytes.NewReader([]byte{0x00})},
			0,
		},
		{
			"0x7F == 0x7F",
			args{bytes.NewReader([]byte{0x7F})},
			127,
		},
		{
			"0x81 0x00 == 0x80",
			args{bytes.NewReader([]byte{0x81, 0x00})},
			128,
		},
		{
			"0xC0 0x00 == 0x2000",
			args{bytes.NewReader([]byte{0xC0, 0x00})},
			8192,
		},
		{
			"0xFF 0x7F == 0x3FFF",
			args{bytes.NewReader([]byte{0xFF, 0x7F})},
			16383,
		},
		{
			"0xFF 0xFF 0xFF 0x7F == 0x0FFFFFFF",
			args{bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0x7F})},
			0xFFFFFFF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := ReadVLQ(tt.args.input); got != tt.want {
				t.Errorf("ParseMidiVLQ() = %v, want %v", got, tt.want)
			}
		})
	}
}
