package befunge

import (
	"bytes"
	"unsafe"

	"github.com/wzshiming/ctc"
)

// Scanner befunge scanner.
type Scanner struct {
	src    [][]byte
	rudder byte
	x, y   int
	ch     byte
	isStr  bool
}

// NewScanner create a new befunge codes scanner.
func NewScanner(s []byte) *Scanner {
	scan := &Scanner{
		src: bytes.Split(s, []byte{'\n'}),
	}
	max := 0
	for _, v := range scan.src {
		if len(v) > max {
			max = len(v)
		}
	}

	bs := make([]byte, max)
	for i, v := range scan.src {
		if off := max - len(v); off >= 0 {
			scan.src[i] = append(scan.src[i], bs[:off]...)
		}
	}
	return scan
}

// String returns the current codes.
func (s *Scanner) String() string {
	str := make([]byte, 0, (1+len(s.src))*(1+len(s.src[0]))+16)
	for i, src := range s.src {
		if i == s.y {
			off := s.x
			str = append(str, src[:off]...)
			str = append(str, ctc.BackgroundRed.String()...)
			str = append(str, src[off])
			str = append(str, ctc.Reset.String()...)
			str = append(str, src[off+1:]...)
		} else {
			str = append(str, src...)
		}
		str = append(str, '\n')
	}
	return *(*string)(unsafe.Pointer(&str))
}

// Scan returns scan a code.
func (s *Scanner) Scan() (int, byte) {
	if str, ok := s.scanString(); ok {
		return str, OpOther
	} else if num, ok := s.scanInteger(); ok {
		return num, OpOther
	} else {
		return 0, s.ch
	}
}

// GetCode get code.
func (s *Scanner) GetCode(x, y int) int {
	if !s.check(int(x), int(y)) {
		return 0
	}
	return int(s.src[y][x])
}

// PutCode put code.
func (s *Scanner) PutCode(x, y, v int) {
	if !s.checkPut(int(x), int(y)) {
		return
	}
	s.src[y][x] = byte(v)
}

func (s *Scanner) scanInteger() (int, bool) {
	if s.ch < '0' || s.ch > '9' {
		return 0, false
	}

	sum := int(s.ch) - '0'
	return int(sum), true
}

func (s *Scanner) scanString() (int, bool) {
	if s.isStr {
		if s.ch == OpStringMode {
			s.isStr = false
			s.Next(1)
			return 0, false
		}
	} else {
		if s.ch != OpStringMode {
			return 0, false
		}
		s.isStr = true
		s.Next(1)

	}
	return int(s.ch), true
}

// SetRudder set rudder.
func (s *Scanner) SetRudder(ru byte) {
	s.rudder = ru
	return
}

// Next scan next
func (s *Scanner) Next(i int) {

	switch s.rudder {
	case 0:
		s.rudder = OpModRight
	case OpMovUp:
		s.y -= i
	case OpMovDown:
		s.y += i
	case OpMovLeft:
		s.x -= i
	case OpModRight:
		s.x += i
	}
	s.checkLoop()

	s.ch = s.src[s.y][s.x]
}

func (s *Scanner) check(x, y int) bool {
	return !(s.x < 0 ||
		s.y < 0 ||
		s.y >= len(s.src) ||
		s.x >= len(s.src[s.y]))
}

func (s *Scanner) checkLoop() {
	if s.y < 0 {
		s.y += len(s.src)
	} else if s.y >= len(s.src) {
		s.y -= len(s.src)
	}
	if s.x < 0 {
		s.x += len(s.src[s.y])
	} else if s.x >= len(s.src[s.y]) {
		s.x -= len(s.src[s.y])
	}
}

func (s *Scanner) checkPut(x, y int) bool {
	if x < 0 ||
		y < 0 {
		return false
	}
	if ya := y - len(s.src); ya > 0 {
		s.src = append(s.src, make([][]byte, ya+2)...)
	}
	if xa := x - len(s.src[y]); xa > 0 {
		s.src[y] = append(s.src[y], make([]byte, xa+1)...)
	}
	return true
}
