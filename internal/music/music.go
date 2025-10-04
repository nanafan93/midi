package music

import (
	"fmt"
)

type keySig struct {
	SF int8 // -7..+7 -ve = number of flats, 0 = C, +ve = number of sharps
	MI uint // 0 = major, 1 = minor
}

type TimeSignature struct {
	Numerator                   uint // top sign
	Denominator                 uint // denominator is power of 2
	MidiClocksPerMetronomeClick byte
	ThirtySecondsPerQuarterNote byte
}

func (ts TimeSignature) String() string {
	return fmt.Sprintf("%d/%d\nClocks/click: %d\n32nd/Quarter note: %d",
		ts.Numerator, ts.Denominator, ts.MidiClocksPerMetronomeClick, ts.ThirtySecondsPerQuarterNote)
}

func GetKeySignature(sf int8, mi uint) string {
	return keySignatures[keySig{SF: sf, MI: mi}]
}

var keySignatures = map[keySig]string{
	{SF: -7, MI: 0}: "C♭ major",
	{SF: -6, MI: 0}: "G♭ major",
	{SF: -5, MI: 0}: "D♭ major",
	{SF: -4, MI: 0}: "A♭ major",
	{SF: -3, MI: 0}: "E♭ major",
	{SF: -2, MI: 0}: "B♭ major",
	{SF: -1, MI: 0}: "F major",
	{SF: 0, MI: 0}:  "C major",
	{SF: 1, MI: 0}:  "G major",
	{SF: 2, MI: 0}:  "D major",
	{SF: 3, MI: 0}:  "A major",
	{SF: 4, MI: 0}:  "E major",
	{SF: 5, MI: 0}:  "B major",
	{SF: 6, MI: 0}:  "F♯ major",
	{SF: 7, MI: 0}:  "C♯ major",

	{SF: -7, MI: 1}: "A♭ minor",
	{SF: -6, MI: 1}: "E♭ minor",
	{SF: -5, MI: 1}: "B♭ minor",
	{SF: -4, MI: 1}: "F minor",
	{SF: -3, MI: 1}: "C minor",
	{SF: -2, MI: 1}: "G minor",
	{SF: -1, MI: 1}: "D minor",
	{SF: 0, MI: 1}:  "A minor",
	{SF: 1, MI: 1}:  "E minor",
	{SF: 2, MI: 1}:  "B minor",
	{SF: 3, MI: 1}:  "F♯ minor",
	{SF: 4, MI: 1}:  "C♯ minor",
	{SF: 5, MI: 1}:  "G♯ minor",
	{SF: 6, MI: 1}:  "D♯ minor",
	{SF: 7, MI: 1}:  "A♯ minor",
}

const middleC = 60
const octaveMiddleC = 4
const semitonesPerOctave = 12

var semitones = map[int]string{
	0:  "C",
	1:  "C♯, D♭",
	2:  "D",
	3:  "D♯, E♭",
	4:  "E",
	5:  "F",
	6:  "F♯, G♭",
	7:  "G",
	8:  "G♯, A♭",
	9:  "A",
	10: "A♯, B♭",
	11: "B",
}

/*
IntegerToNoteName converts a MIDI note number (0-127) to a human-readable note name with octave number.
The octave number follows the scientific pitch notation where Middle C (MIDI note 60) is C4.
*/
func IntegerToNoteName(noteNumber int) string {
	distanceFromMiddleC := noteNumber - middleC
	octaveNumber := noteNumber/12 - 1
	mod := distanceFromMiddleC % semitonesPerOctave
	if mod < 0 {
		mod += semitonesPerOctave
	}
	noteNameWithOctaveNumber := fmt.Sprintf("(%s)%d", semitones[mod], octaveNumber)
	return noteNameWithOctaveNumber
}
