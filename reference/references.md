# Useful references

## Important Structures
### MIDI structure
```
MIDI file is made up of chunks like so: 
MThd <length of header data>
<header data>
MTrk <length of track data>
<track data>
MTrk <length of track data>
<track data>
...
EOF
```
### Header chunk
```
MThd <length of header data>
<header data>
```
* Mthd: 4 bytes ASCII "MThd"
* length of header data: 4 bytes big-endian uint32, always 6
* header data: 6 bytes
  * format type: 2 bytes big-endian uint16, 0, 1 or 2
  * number of tracks: 2 bytes big-endian uint16
  * time division: 2 bytes big-endian uint16
    * if the most significant bit is 0, the remaining 15 bits represent ticks per quarter note
    * if the most significant bit is 1, the remaining 15 bits represent frames per second (negative) and ticks per frame

### Track chunk
```
MTrk <length of track data>
<track data>

where track data is a sequence of MIDI events:

<Track Chunk> = <chunk type> <length> <MTrk event>+
<Track Event> = <delta time> <event>
<event> = <MIDI event> | <Meta event> | <SysEx event>
<MIDI event> = <status byte> <data byte 1> [<data byte 2>]
<Meta event> = 0xFF <meta type> <length> <data>
<SysEx event> = 0xF0 <length> <data> | 0xF7 <length> <data> 

```
* MTrk: 4 bytes ASCII "MTrk"
* length of track data: 4 bytes big-endian uint32, length of the track
* track data: sequence of MIDI events, each event has:
  * delta time: VLQ (Variable Length Quantity)
  * event type and parameters: varies based on event type
  * MIDI events can be:
    * MIDI channel messages (e.g., Note On, Note Off, Control Change)
    * Meta events (e.g., Tempo Change, Time Signature)
    * System Exclusive events


## VLQ (Variable Length Quantity)
- [Variable Length Quantity Wiki][1]
- [Codegolf vlq][2]


[1]: https://en.wikipedia.org/wiki/Variable-length_quantity "Variable Length Quantity Wiki"
[2]: https://codegolf.stackexchange.com/questions/189471/decode-a-variable-length-quantity "Codegolf vlq"