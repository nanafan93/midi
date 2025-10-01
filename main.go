package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type MIDIHeader struct {
	ChunkType [4]byte // "MThd"
	Length    uint32  // always 6
	Format    uint16
	NumTracks uint16
	Division  uint16
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
	header := make([]byte, 14)
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	printEntireFile(file)
	defer file.Close()
	file.Seek(0, io.SeekStart)
	n, err := file.Read(header)
	if err != nil || n < 14 {
		log.Fatal("Cannot read header: ", err)
	}
	if string(header[:4]) != "MThd" {
		log.Fatal("Not a valid MIDI file (missing MThd header)")
	}
	length := binary.BigEndian.Uint32(header[4:8])
	format := binary.BigEndian.Uint16(header[8:10])
	nTracks := binary.BigEndian.Uint16(header[10:12])
	division := binary.BigEndian.Uint16(header[12:14])
	fmt.Printf("Header length: %d\nFormat: %d\nTracks: %d\nDivision: %d\n", length, format, nTracks, division)

}
