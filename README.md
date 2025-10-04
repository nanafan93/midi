# Learn Standard MIDI File Format
Parse and display information in a Standard MIDI File (SMF) using Go.
## Run
Use a MIDI file as an argument to run the program, for example:
```bash
go run main.go assets/mary_had_a_little_lamb.mid
```
## Learning outcomes
- Learn how to parse binary file formats in Go
- Practice working with file I/O and byte manipulation in Go
- Understand VLQ (Variable Length Quantity) encoding used in MIDI files

## Improvements
- Improved error handling
- Playing a track ?!??
- Add tests lel
- Implement a state machine for parsing