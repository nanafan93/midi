package midi

import (
	"fmt"
	"midi/internal/music"
)

type FormatType int

type Track struct {
	ChunkType [4]byte
	Length    uint32
}

const (
	SingleMultiChannelTrack FormatType = iota
	MultiTrack
	SequentialIndpendentTracks
)

func (h Header) String() string {
	isTicksPerQuarterNote := h.Division&(1<<15) == 0
	var divisionType string
	if isTicksPerQuarterNote {
		divisionType = "Ticks Per Quarter Note"
	} else {
		divisionType = "SMPTE + MIDI Code"
	}
	return fmt.Sprintf("Chunk Type: %s\nLength: %d\nFormat: %s\nNumber of tracks: %d\n%s: %d\n",
		string(h.ChunkType[:]), h.Length, h.getFormat(), h.NumTracks, divisionType, h.Division)
}

func (h Header) GetChunkType() string {
	return string(h.ChunkType[:])
}

type Header struct {
	ChunkType [4]byte // "MThd"
	Length    uint32  // always 6
	Format    uint16
	NumTracks uint16
	Division  uint16
}

func (h Header) getFormat() FormatType {
	switch h.Format {
	case 0:
		return SingleMultiChannelTrack
	case 1:
		return MultiTrack
	case 2:
		return SequentialIndpendentTracks
	default:
		panic("Invalid format type")
	}
}

func (f FormatType) String() (s string) {
	switch f {
	case MultiTrack:
		s = fmt.Sprintf("%d: Multi Track", f)
	case SingleMultiChannelTrack:
		s = fmt.Sprintf("%d: Single Track", f)
	case SequentialIndpendentTracks:
		s = fmt.Sprintf("%d: Sequential Independent Tracks", f)
	}
	return
}

type ChannelVoiceMessage struct {
	Name        string
	Description string
	Length      uint
	Decode      func([]byte) string
}

var ChannelVoiceMessages = map[byte]ChannelVoiceMessage{
	0x80: {"Note Off", "Note released", 2, func(data []byte) string {
		return fmt.Sprintf("key=%s velocity=%d", music.IntegerToNoteName(int(data[0])), data[1])
	}},
	0x90: {"Note On", "Note pressed", 2, func(data []byte) string {
		return fmt.Sprintf("key=%s velocity=%d", music.IntegerToNoteName(int(data[0])), data[1])
	}},
	0xA0: {"Poly Aftertouch", "Pressure on a key", 2, func(data []byte) string {
		return fmt.Sprintf("key=%d pressure=%d", data[0], data[1])
	}},
	0xB0: {"Control Change", "Controller change", 2, func(data []byte) string {
		return fmt.Sprintf("controller=%d value=%d", data[0], data[1])
	}},
	0xC0: {"Program Change", "Instrument change", 1, func(data []byte) string {
		return fmt.Sprintf("program=%d", data[0])
	}},
	0xD0: {"Channel Pressure", "Pressure on channel", 1, func(data []byte) string {
		return fmt.Sprintf("pressure=%d", data[0])
	}},
	0xE0: {"Pitch Bend", "Pitch wheel change", 2, func(data []byte) string {
		value := int(data[0]) | (int(data[1]) << 7) // 14-bit
		return fmt.Sprintf("bend=%d", value-8192)   // center = 0
	}},
}

type MetaEvent struct {
	Name        string
	Description string
	FixedLength int // -1 if variable length
	Decode      func([]byte) string
}

var MetaEvents = map[byte]MetaEvent{
	0x21: {
		"Midi Port",
		"Unoffical extension MIDI Port number",
		1,
		func(b []byte) string {
			return fmt.Sprintf("Port %d", b[0])
		},
	},
	0x2F: {
		"End of track",
		"Not optional event indicating end of track",
		0,
		func(bytes []byte) string {
			return "End of Track"
		},
	},
	0x01: {
		"Text Event",
		"Any amount of text describing anything",
		-1,
		func(bytes []byte) string {
			return string(bytes)
		},
	},
	0x02: {
		"Copyright notice",
		"Copyright notice in ASCII printable text",
		-1,
		func(bytes []byte) string {
			return string(bytes)
		},
	},
	0x03: {
		"Sequence/Track Name",
		"If in a format 0 track, or the first track in a format 1 file, the name of the sequence.\nOtherwise, the name of the track.",
		-1,
		func(bytes []byte) string {
			return string(bytes)
		}},
	0x51: {
		"Set Tempo",
		"Tempo in microseconds per quarter note",
		3,
		func(data []byte) string {
			usPerQuarter := int(data[0])<<16 | int(data[1])<<8 | int(data[2])
			bpm := 60000000 / usPerQuarter
			return fmt.Sprintf("%d μs/qn (≈ %d BPM)", usPerQuarter, bpm)
		},
	},
	0x58: {"Time Signature",
		"Numerator, denominator, clocks, 32nd notes",
		4,
		func(data []byte) string {
			return music.TimeSignature{
				Numerator:                   uint(data[0]),
				Denominator:                 1 << data[1],
				MidiClocksPerMetronomeClick: data[2],
				ThirtySecondsPerQuarterNote: data[3],
			}.String()
		},
	},
	0x59: {"Key Signature", "Specifies key (sf, mi) ", 2,
		func(data []byte) string {
			sf := int8(data[0])
			mi := uint(data[1])
			return music.GetKeySignature(sf, mi)
		},
	},
}
