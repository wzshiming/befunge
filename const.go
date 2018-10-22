package befunge

// https://en.wikipedia.org/wiki/Befunge
// 0-9	Push this number on the stack.
const (
	OpAdd         = '+' // Addition: Pop a and b, then push a+b.
	OpSub         = '-' // Subtraction: Pop a and b, then push b-a.
	OpMult        = '*' // Multiplication: Pop a and b, then push a*b.
	OpDiv         = '/' // Integer division: Pop a and b, then push b/a, rounded down. If a is zero, push zero.
	OpMod         = '%' // Modulo: Pop a and b, then push the b%a. If a is zero, push zero.
	OpNot         = '!' // Logical NOT: Pop a value. If the value is zero, push 1; otherwise, push zero.
	OpGreaterThan = '`' // Greater than: Pop a and b, then push 1 if b>a, otherwise push zero.

	OpDup        = ':'  // Duplicate value on top of the stack. If there is nothing on top of the stack, push a 0.
	OpSwap       = '\\' // Swap two values on top of the stack. If there is only one value, pretend there is an extra 0 on bottom of the stack.
	OpStringMode = '"'  // Start string mode: push each character's ASCII value all the way up to the next ".
	OpPop        = '$'  // Pop value from the stack and discard it.

	OpOutInt  = '.' // Pop value and output as an integer.
	OpOutRune = ',' // Pop value and output the ASCII character represented by the integer code that is stored in the value.
	OpInInt   = '&' // Ask user for a number and push it
	OpInRune  = '~' // Ask user for a character and push its ASCII value

	OpGetCode = 'g' // A "get" call (a way to retrieve data in storage). Pop y and x, then push ASCII value of the character at that position in the program
	OpPutCode = 'p' // A "put" call (a way to store a value for later use). Pop y, x and v, then change the character at the position (x,y) in the program to the character with ASCII value v.

	OpModRight  = '>' // Start moving right.
	OpMovLeft   = '<' // Start moving left.
	OpMovUp     = '^' // Start moving up.
	OpMovDown   = 'v' // Start moving down.
	OpMovRandom = '?' // Start moving in a random cardinal direction.
	OpIfHoriz   = '_' // Pop a value; move right if value = 0, left otherwise.
	OpIfVert    = '|' // Pop a value; move down if value = 0, up otherwise.
	OpBridge    = '#' // Trampoline: Skip next cell.
	OpBlank     = ' ' // (i.e. a space) No-op. Does nothing.
	OpEnd       = '@' // End program.
)
