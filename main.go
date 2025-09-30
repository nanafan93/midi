package main

import (
	"fmt"
	"log"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: midi <midi-file>")
		os.Exit(1)
	}
	fileName := os.Args[1]
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}
	fmt.Printf("MIDI size: %d bytes\n", len(fileBytes))
	for i, fileByte := range fileBytes {
		fmt.Printf("%02X ", fileByte)
		if (i+1)%10 == 0 {
			fmt.Println()
		}
	}
}
