package class

import (
	"errors"
	"fmt"
)

var (
	ErrPartialCharacter = errors.New("malformed input: partial character at end")
)

func DecodeRunes(b []byte) (res []rune, err error) {
	count := 0
	cpCount := 0
	c, c2, c3 := 0, 0, 0
	utfLen := len(b)
	cps := make([]int, utfLen)
	for count < utfLen {
		c = int(b[count]) & 0xff
		switch c >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			// 0xxxxxxx
			count++
			cps[cpCount] = c
			cpCount++
		case 12, 13:
			// 110x xxxx   10xx xxxx
			count += 2
			if count > utfLen {
				err = ErrPartialCharacter
				return
			}
			c2 = int(b[count-1])
			if (c2 & 0xC0) != 0x80 {
				err = fmt.Errorf("malformed input around byte %d", count)
				return
			}
			cps[cpCount] = int(((c & 0x1F) << 6) | (c2 & 0x3F))
			cpCount++
		case 14:
			// 1110 xxxx  10xx xxxx  10xx xxxx
			count += 3
			if count > utfLen {
				err = ErrPartialCharacter
				return
			}
			c2 = int(b[count-2])
			c3 = int(b[count-1])
			if ((c2 & 0xC0) != 0x80) || ((c3 & 0xC0) != 0x80) {
				err = fmt.Errorf("malformed input around byte %d", count-1)
				return
			}
			cps[cpCount] = int(((c & 0x0F) << 12) | ((c2 & 0x3F) << 6) | ((c3 & 0x3F) << 0))
			cpCount++
		default:
			err = fmt.Errorf("malformed input around byte %d", count)
			return
		}
	}
	res = make([]rune, len(cps))
	for i := 0; i < len(cps); i++ {
		res[i] = rune(cps[i])
	}
	return
}
