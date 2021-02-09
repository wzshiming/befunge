package befunge

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/wzshiming/ctc"
)

var (
	pointColorStart   = ctc.BackgroundRed.String()
	editedColorStart  = ctc.BackgroundBlue.String()
	colorEnd          = ctc.Reset.String()
	editedNoPrintText = fmt.Sprintf("%s?%s", ctc.BackgroundBlue|ctc.ForegroundRed, colorEnd)
)

// Scanner befunge scanner.
type Scanner struct {
	src           [][]byte
	edited        map[int]map[int]struct{}
	rudder        byte
	x, y          int
	ch            byte
	isStr         bool
	height, width int
}

// NewScanner create a new befunge codes scanner.
func NewScanner(src []byte) *Scanner {
	s := &Scanner{
		src: bytes.Split(src, []byte{'\n'}),
	}
	for _, v := range s.src {
		if len(v) > s.width {
			s.width = len(v)
		}
	}
	s.height = len(s.src)
	s.resize(s.width, s.height)
	return s
}

// String returns the current codes.
func (s *Scanner) String() string {
	str := make([]byte, 0, (1+len(s.src))*(1+len(s.src[0]))+16)
	for i, src := range s.src {
		if i == s.y {
			off := s.x
			str = append(str, s.markEdited(0, i, src[:off])...)
			str = append(str, pointColorStart...)
			str = append(str, src[off])
			str = append(str, colorEnd...)
			str = append(str, s.markEdited(off+1, i, src[off+1:])...)
		} else {
			str = append(str, s.markEdited(0, i, src)...)
		}
		str = append(str, '\n')
	}
	str = bytes.TrimSpace(str)
	str = append(str, []byte(fmt.Sprintf("\n(%d,%d) %s", s.x, s.y, CodeText(int(s.ch))))...)
	return *(*string)(unsafe.Pointer(&str))
}

func (s *Scanner) markEdited(xoff, y int, data []byte) []byte {
	if s.edited == nil || s.edited[y] == nil {
		return data
	}
	out := make([]byte, 0, len(data)*2)
	for i, c := range data {
		_, edited := s.edited[y][i+xoff]
		if !edited {
			out = append(out, c)
		} else if strconv.IsPrint(rune(c)) {
			out = append(out, editedColorStart...)
			out = append(out, c)
			out = append(out, colorEnd...)
		} else {
			out = append(out, editedNoPrintText...)
		}
	}
	return out
}

// SprintText returns the print text
func SprintText(r []byte) []byte {
	out := []byte{}
	for _, c := range r {
		if strconv.IsPrint(rune(c)) {
			out = append(out, c)
		} else if c == 0 {
			out = append(out, ' ')
		} else {
			out = append(out, editedNoPrintText...)
		}
	}
	return out
}

// Point returns current point
func (s *Scanner) Point() (int, int) {
	return s.x, s.y
}

// Size returns current size
func (s *Scanner) Size() (int, int) {
	return s.width, s.height
}

// Scan returns scan a code.
func (s *Scanner) Scan() (int, byte) {
	if str, op, ok := s.scanString(); ok {
		return str, op
	} else if num, op, ok := s.scanInteger(); ok {
		return num, op
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
	if s.edited == nil {
		s.edited = map[int]map[int]struct{}{}
	}
	if s.edited[y] == nil {
		s.edited[y] = map[int]struct{}{}
	}
	s.edited[y][x] = struct{}{}
	s.src[y][x] = byte(v)
}

func (s *Scanner) scanInteger() (int, byte, bool) {
	if s.ch < '0' || s.ch > '9' {
		return 0, OpNone, false
	}

	sum := int(s.ch) - '0'
	return int(sum), OpOther, true
}

func (s *Scanner) scanString() (int, byte, bool) {
	if s.isStr {
		if s.ch == OpStringMode {
			s.isStr = false
			return 0, OpNone, true
		}
	} else {
		if s.ch == OpStringMode {
			s.isStr = true
			return 0, OpNone, true
		}
		return 0, OpNone, false
	}
	if s.ch == 0 {
		s.ch = ' '
	}
	return int(s.ch), OpOther, true
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
	return !(x < 0 ||
		y < 0 ||
		y >= len(s.src) ||
		x >= len(s.src[y]))
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
	if x >= s.width || y >= s.height {
		if x >= s.width {
			s.width = x
		}
		if y >= s.height {
			s.height = y
		}
		s.resize(s.width, s.height)
	}
	return true
}

func (s *Scanner) resize(x, y int) {
	if off := y - len(s.src); off >= 0 {
		s.src = append(s.src, make([][]byte, off+1)...)
	}
	for i := 0; i != len(s.src); i++ {
		v := s.src[i]
		if off := x - len(v); off >= 0 {
			s.src[i] = append(s.src[i], make([]byte, off+1)...)
		}
	}
}

// CodeText returns the text of code
func CodeText(v int) string {
	return fmt.Sprintf("(%q %d)", v, v)
}
