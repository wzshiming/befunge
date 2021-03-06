package befunge

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
)

// Runner befunge runner.
type Runner struct {
	Scanner
	stack    []int
	input    io.Reader
	output   io.Writer
	step     func()
	debug    bool
	errors   []error
	randFunc func(n int) int
}

// NewRunner create a new befunge codes runner.
func NewRunner(s []byte) *Runner {
	return &Runner{
		Scanner:  *NewScanner(s),
		output:   os.Stdout,
		input:    os.Stdin,
		randFunc: rand.Intn,
	}
}

// Stack returns stack.
func (r *Runner) Stack() []int {
	return r.stack
}

// SetInput sets input.
func (r *Runner) SetInput(input io.Reader) {
	r.input = input
}

// SetOutput sets output.
func (r *Runner) SetOutput(output io.Writer) {
	r.output = output
}

// SetRandFunc sets randFunc.
func (r *Runner) SetRandFunc(randFunc func(n int) int) {
	r.randFunc = randFunc
}

// SetStep sets step.
func (r *Runner) SetStep(f func()) {
	r.step = f
}

// SetDebug sets debug.
func (r *Runner) SetDebug(e bool) {
	r.debug = e
}

// Put value to the stack.
func (r *Runner) Put(i int) {
	r.stack = append(r.stack, i)
}

// Get value from of the stack.
func (r *Runner) Get() int {
	if len(r.stack) == 0 {
		return 0
	}
	v := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return v
}

// Swap two values on top of the stack.
func (r *Runner) Swap() {
	switch len(r.stack) {
	case 0:
		return
	case 1:
		r.stack = append(r.stack, 0)
	}
	r.stack[len(r.stack)-1], r.stack[len(r.stack)-2] = r.stack[len(r.stack)-2], r.stack[len(r.stack)-1]
}

// Duplicate value on top of the stack.
func (r *Runner) Duplicate() {
	if len(r.stack) == 0 {
		r.stack = append(r.stack, 0)
		return
	}
	r.stack = append(r.stack, r.stack[len(r.stack)-1])
	return
}

func (r *Runner) runStep() (bool, error) {
	val, ch := r.Scan()

	switch ch {
	case OpOther:
		r.Put(val)
	case OpAdd:
		a := r.Get()
		b := r.Get()
		r.Put(a + b)
	case OpSub:
		a := r.Get()
		b := r.Get()
		r.Put(b - a)
	case OpMult:
		a := r.Get()
		b := r.Get()
		r.Put(a * b)
	case OpDiv:
		a := r.Get()
		b := r.Get()
		r.Put(b / a)
	case OpMod:
		a := r.Get()
		b := r.Get()
		r.Put(b % a)
	case OpNot:
		if r.Get() == 0 {
			r.Put(1)
		} else {
			r.Put(0)
		}
	case OpGreaterThan:
		a := r.Get()
		b := r.Get()
		if b > a {
			r.Put(1)
		} else {
			r.Put(0)
		}
	case OpIfHoriz:
		if r.Get() == 0 {
			r.SetRudder(OpModRight)
		} else {
			r.SetRudder(OpMovLeft)
		}
	case OpIfVert:
		if r.Get() == 0 {
			r.SetRudder(OpMovDown)
		} else {
			r.SetRudder(OpMovUp)
		}
	case OpDup:
		r.Duplicate()
	case OpSwap:
		r.Swap()
	case OpPop:
		r.Get()
	case OpPutCode:
		y := r.Get()
		x := r.Get()
		v := r.Get()
		r.PutCode(x, y, v)
	case OpGetCode:
		y := r.Get()
		x := r.Get()
		v := r.GetCode(x, y)
		r.Put(v)
	case OpOutInt:
		r.Output(strconv.FormatInt(int64(r.Get()), 10))
	case OpOutRune:
		r.Output(string([]byte{byte(r.Get())}))
	case OpInInt:
		v := 0
		info := fmt.Sprintf("\n(Enter a number '%s'): ", string([]byte{ch}))
		for {
			r.Output(info)
			_, err := fmt.Fscanf(r.input, "%d\n", &v)
			if err == nil {
				break
			}
			r.errors = append(r.errors, err)
		}
		r.Put(v)
	case OpInRune:
		char := 0
		info := fmt.Sprintf("\n(Enter a character '%s'): ", string([]byte{ch}))
		for {
			r.Output(info)
			_, err := fmt.Fscanf(r.input, "%c\n", &char)
			if err == nil {
				break
			}
			r.errors = append(r.errors, err)
		}
		r.Put(char)
	case OpModRight, OpMovLeft, OpMovUp, OpMovDown:
		r.SetRudder(ch)
	case OpBridge:
		r.Next(1)
	case OpMovRandom:
		randSwitch := []byte{OpModRight, OpMovLeft, OpMovUp, OpMovDown}
		ru := randSwitch[r.randFunc(len(randSwitch))]
		r.SetRudder(ru)
	case OpEnd:
		return false, nil

	case OpBlank, OpNone:
	default:
		x, y := r.Point()
		err := fmt.Errorf("Error in %d, %d undefined: %s", x, y, CodeText(int(ch)))
		r.errors = append(r.errors, err)
	}
	return true, nil
}

// Output print
func (r *Runner) Output(s string) {
	if r.output != nil {
		io.WriteString(r.output, s)
	}
	if r.step != nil {
		r.step()
	}
}

// Errors of the Runner
func (r *Runner) Errors() []error {
	return r.errors
}

// Run befunge code.
func (r *Runner) Run() error {
	for {
		r.Next(1)
		if r.step != nil {
			r.step()
		}
		ok, err := r.runStep()
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}
