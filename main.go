package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"midi/internal/midi"
	"midi/internal/vlq"
	"os"
)

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
	r := bufio.NewReader(file)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	printEntireFile(file)
	defer file.Close()
	file.Seek(0, io.SeekStart)
	var midiHeader midi.Header
	err = binary.Read(r, binary.BigEndian, &midiHeader)
	if err != nil {
		log.Fatal("Cant use binary.Read")
	}
	if string(midiHeader.ChunkType[:]) != "MThd" {
		log.Fatal("Bad MIDI header")
	}
	fmt.Printf("Header %+v", midiHeader)
	numTracksFound := 0
	//trackMap := make(map[int][]byte)
	for {
		var trackChunk midi.Track

		err = binary.Read(r, binary.BigEndian, &trackChunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Error reading track")
		}
		if string(trackChunk.ChunkType[:]) != "MTrk" {
			log.Fatal("Not a track")
		}
		trackData := make([]byte, trackChunk.Length)
		numTracksFound++
		fmt.Printf("Encountered track with length %d\n", trackChunk.Length)
		_, err = io.ReadFull(r, trackData)
		if err != nil {
			log.Fatal("Error reading track data")
		}
		printTrackEvents(trackData)
		//trackMap[numTracksFound] = trackData

	}
	fmt.Printf("Found %d tracks", numTracksFound)
	//fmt.Printf("Track map %v", trackMap)

}

func printMetaEvent(r *bytes.Reader) {
	metaEventType, _ := r.ReadByte()
	metaEvent, exists := midi.MetaEvents[metaEventType]
	if !exists {
		log.Fatalf("Meta event %02X is not yet supported", metaEventType)
	}
	var b []byte
	if metaEvent.FixedLength != -1 {
		b = make([]byte, metaEvent.FixedLength)
		r.ReadByte()
	} else {
		metaEventLength, _ := vlq.ReadVLQ(r)
		b = make([]byte, metaEventLength)
	}
	if r.Len() == 0 {
		return
	}
	_, err := r.Read(b)
	if err != nil {
		log.Fatal("Can't read meta event data")
	}
	fmt.Printf("%s: %s\n", metaEvent.Name, metaEvent.Decode(b))
}
func printTrackEvents(trackData []byte) {
	reader := bytes.NewReader(trackData)
	lastStatus := byte(0)

	for reader.Len() > 0 {
		// Always start by reading delta-time
		deltaTime, err := vlq.ReadVLQ(reader)
		if err != nil {
			log.Fatal("Error reading event")
		}
		fmt.Printf("Delta time %d\n", deltaTime)

		// Peek at next byte (could be a status or data for running status)
		next, err := reader.ReadByte()
		if err != nil {
			log.Fatal("Error reading next byte")
		}

		var statusByte byte
		if next >= 0x80 {
			// This is a real status byte
			statusByte = next
			lastStatus = statusByte
		} else {
			// Running status: reuse last status byte
			if lastStatus == 0 {
				log.Fatal("attempt to use running status before setting")
			}
			statusByte = lastStatus
			// Put back the data byte we just consumed
			reader.UnreadByte()
		}

		// Dispatch by type
		if statusByte == 0xFF {
			printMetaEvent(reader)
		} else if statusByte == 0xF0 || statusByte == 0xF7 {
			skipEvent(reader, statusByte)
		} else if statusByte&0xF0 >= 0x80 && statusByte&0xF0 <= 0xE0 {
			printMidiChannelEvent(reader, statusByte)
		} else {
			printRunningStatusEvent(reader, lastStatus)
		}
	}
}

func printRunningStatusEvent(reader *bytes.Reader, status byte) {
	channelMessage := midi.ChannelVoiceMessages[status]
	b := make([]byte, channelMessage.Length)
	_, err := reader.Read(b)
	if err != nil {
		log.Fatal("Can't read meta event data")
	}

	fmt.Printf("%s: %s\n", channelMessage.Name, channelMessage.Decode(b))
}

func printMidiChannelEvent(reader *bytes.Reader, statusByte byte) {
	highNibbleKey := statusByte & 0xF0
	metaEvent, exists := midi.ChannelVoiceMessages[highNibbleKey]
	if !exists {
		log.Fatalf("Channel message %02X is not yet supported", statusByte)
	}
	b := make([]byte, metaEvent.Length)
	_, err := reader.Read(b)
	if err != nil {
		log.Fatal("Can't read meta event data")
	}
	fmt.Printf("%s: %s\n", metaEvent.Name, metaEvent.Decode(b))
}

// skipEvent consumes bytes from a reader to move past SysEx or MIDI events
func skipEvent(r *bytes.Reader, statusByte byte) error {
	switch {
	case statusByte == 0xF0 || statusByte == 0xF7: // SysEx
		// SysEx events: <0xF0 or 0xF7> <VLQ length> <data>
		length, err := vlq.ReadVLQ(r)
		if err != nil {
			return err
		}
		_, err = r.Seek(int64(length), io.SeekCurrent)
		return err

	case statusByte >= 0x80 && statusByte <= 0xEF: // MIDI Channel Event
		// High nibble = event type, low nibble = channel
		eventType := statusByte & 0xF0
		var dataBytes int
		switch eventType {
		case 0xC0, 0xD0: // Program Change, Channel Pressure
			dataBytes = 1
		default: // Note On/Off, Control Change, Pitch Bend, etc.
			dataBytes = 2
		}
		_, err := r.Seek(int64(dataBytes), io.SeekCurrent)
		return err

	default:
		return fmt.Errorf("unknown status byte: 0x%X", statusByte)
	}
}
