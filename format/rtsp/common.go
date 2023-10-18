package rtsp

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

const charset = "0123456789abcdefghijklmnopqrstuvwxyz"

func genSessionCookie() string {
	const idLength = 23
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	idBytes := make([]byte, idLength)

	for i := range idBytes {
		idBytes[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(idBytes)
}

func findRTSP(c *Client) (block []byte, data []byte, err error) {
	const (
		R = iota + 1
		T
		S
		Header
		Dollar
	)

	var _peek [8]byte
	peek := _peek[0:0]
	stat := 0

	for i := 0; ; i++ {
		var b byte
		if b, err = c.brconn.ReadByte(); err != nil {
			return
		}
		switch b {
		case 'R':
			if stat == 0 {
				stat = R
			}
		case 'T':
			if stat == R {
				stat = T
			}
		case 'S':
			if stat == T {
				stat = S
			}
		case 'P':
			if stat == S {
				stat = Header
			}
		case '$':
			if stat != Dollar {
				stat = Dollar
				peek = _peek[0:0]
			}
		default:
			if stat != Dollar {
				stat = 0
				peek = _peek[0:0]
			}
		}

		if false && c.DebugRtp {
			fmt.Println("rtsp: findRTSP", i, b)
		}

		if stat != 0 {
			peek = append(peek, b)
		}
		if stat == Header {
			data = peek
			return
		}

		if stat == Dollar && len(peek) >= 12 {
			if c.DebugRtp {
				fmt.Println("rtsp: dollar at", i, len(peek))
			}
			if blocklen, _, ok := c.parseBlockHeader(peek); ok {
				left := blocklen + 4 - len(peek)
				if left >= 0 {
					block = append(peek, make([]byte, left)...)
					if _, err = io.ReadFull(c.brconn, block[len(peek):]); err != nil {
						return
					}
					return
				} else {
					fmt.Println("Left < 0 ", blocklen, len(peek), left)
				}
			}
			stat = 0
			peek = _peek[0:0]
		}
	}

	return
}

func findHTTP(c *Client) (block []byte, data []byte, err error) {
	const (
		H = iota + 1
		T1
		T2
		Header
		Dollar
	)

	var _peek [8]byte
	peek := _peek[0:0]
	stat := 0

	for i := 0; ; i++ {
		var b byte
		if b, err = c.brconn.ReadByte(); err != nil {
			return
		}
		switch b {
		case 'H':
			if stat == 0 {
				stat = H
			}
		case 'T':
			if stat == H {
				stat = T1
			} else if stat == T1 {
				stat = T2
			}
		case 'P':
			if stat == T2 {
				stat = Header
			}
		case '$':
			if stat != Dollar {
				stat = Dollar
				peek = _peek[0:0]
			}
		default:
			if stat != Dollar {
				stat = 0
				peek = _peek[0:0]
			}
		}

		if false && c.DebugRtp {
			fmt.Println("rtsp: findRTSP", i, b)
		}

		if stat != 0 {
			peek = append(peek, b)
		}
		if stat == Header {
			data = peek
			return
		}

		if stat == Dollar && len(peek) >= 12 {
			if c.DebugRtp {
				fmt.Println("rtsp: dollar at", i, len(peek))
			}
			if blocklen, _, ok := c.parseBlockHeader(peek); ok {
				left := blocklen + 4 - len(peek)
				if left >= 0 {
					block = append(peek, make([]byte, left)...)
					if _, err = io.ReadFull(c.brconn, block[len(peek):]); err != nil {
						return
					}
					return
				} else {
					fmt.Println("Left < 0 ", blocklen, len(peek), left)
				}
			}
			stat = 0
			peek = _peek[0:0]
		}
	}

	return
}

func readLFLF(c *Client) (block []byte, data []byte, err error) {
	const (
		LF = iota + 1
		LFLF
	)
	peek := []byte{}
	stat := 0
	dollarpos := -1
	lpos := 0
	pos := 0

	for {
		var b byte
		if b, err = c.brconn.ReadByte(); err != nil {
			return
		}
		switch b {
		case '\n':
			if stat == 0 {
				stat = LF
				lpos = pos
			} else if stat == LF {
				if pos-lpos <= 2 {
					stat = LFLF
				} else {
					lpos = pos
				}
			}
		case '$':
			dollarpos = pos
		}
		peek = append(peek, b)

		if stat == LFLF {
			data = peek
			return
		} else if dollarpos != -1 && dollarpos-pos >= 12 {
			hdrlen := dollarpos - pos
			start := len(peek) - hdrlen
			if blocklen, _, ok := c.parseBlockHeader(peek[start:]); ok {
				block = append(peek[start:], make([]byte, blocklen+4-hdrlen)...)
				if _, err = io.ReadFull(c.brconn, block[hdrlen:]); err != nil {
					return
				}
				return
			}
			dollarpos = -1
		}

		pos++
	}

	return
}
