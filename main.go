package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type FormatType int

const (
	SingleMultiChannelTrack FormatType = iota
	MultiTrack
	SequentialIndpendentTracks
)

func (f FormatType) String() (s string) {
	switch f {
	case MultiTrack:
		s = "Multi Track"
	case SingleMultiChannelTrack:
		s = "Single Track"
	case SequentialIndpendentTracks:
		s = "Sequential Independent Tracks"
	}
	return
}

type MidiFile struct {
	MidiHeader MidiHeader
}

type MidiHeader struct {
	ChunkType      [4]byte
	Length         uint32
	Format         uint16
	NumberOfTracks uint16
	Division       uint16
}

func (h MIDIHeader) String() string {
	isTicksPerQuarterNote := h.Division&(1<<15) == 0
	return fmt.Sprintf("Chunk Type: %s\nLength: %d\nFormat: %s\nNumber of tracks: %d\nDivision: %d Ticker Per QN: %t\n",
		string(h.ChunkType[:]), h.Length, h.getFormat(), h.NumTracks, h.Division, isTicksPerQuarterNote)
}

func (h MIDIHeader) getFormat() FormatType {
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

type MIDIHeader struct {
	ChunkType [4]byte // "MThd"
	Length    uint32  // always 6
	Format    uint16
	NumTracks uint16
	Division  uint16
}

func (h MIDIHeader) GetChunkType() string {
	return string(h.ChunkType[:])
}

func printEntireFile(file *os.File) {
	buf := make([]byte, 512)
	totalSize := 0
	printedBytes := 0
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("something went wrong reading the file")
			}
		}
		totalSize += n
		for _, fileByte := range buf[:n] {
			fmt.Printf("%02X ", fileByte)
			printedBytes++
			if (printedBytes+1)%10 == 0 {
				fmt.Println()
			}
		}
	}
	fmt.Printf("\nMIDI size: %d bytes\n", totalSize)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: midi <midi-file>")
		os.Exit(1)
	}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	printEntireFile(file)
	defer file.Close()
	file.Seek(0, io.SeekStart)
	var midiHeader MIDIHeader
	err = binary.Read(file, binary.BigEndian, &midiHeader)
	if err != nil {
		log.Fatal("Cant use binary.Read")
	}
	if string(midiHeader.ChunkType[:]) != "MThd" {
		log.Fatal("Bad MIDI header")
	}
	fmt.Printf("Header %+v", midiHeader)

}
