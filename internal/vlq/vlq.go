package vlq

import "io"

func ReadVLQ(r io.ByteReader) (uint, error) {
	var res uint
	for {
		vlqByte, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		res = (res << 7) | uint(vlqByte&0x7F)
		if vlqByte&0x80 == 0 {
			break
		}
	}
	return res, nil
}
