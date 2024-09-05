package goutil

import "fmt"

type StepBytesMethod struct {
	buf *[]byte
	i   *int
	b   *bool
}

// StepBytes runs a for loop over []byte in a safer way.
//
// This method is intended to read through a []byte, one byte at a time,
// and replace characters at a specific point.
// This method can be useful for building a small compiler.
//
// @cb:
//   - return true to continue loop
//   - return false to break the loop
//
// @*i: the current index of the loop
//
// @b: returns the currend byte relative to the index
//  example:
//   b(0) == b[i]
//   b(1) == b[i+1]
//   b(-1) == b[i-1]
//
// @m: other useful methods
func StepBytes(buf *[]byte, cb func(i *int, b func(int) byte, m StepBytesMethod) bool) {
	end := false
	m := StepBytesMethod{
		buf: buf,
		b:   &end,
	}

	for i := 0; i < len(*buf); i++ {
		m.i = &i

		if !cb(&i, func(s int) byte {
			if i+s < len(*buf) {
				return (*buf)[i+s]
			}
			return 0
		}, m) {
			break
		}

		if end {
			break
		}
	}
}

// Inc increments i and will break the loop if i >= len(buf)
//
// this method also returns false, if the loop should break,
// and true if the loop should continue
func (m *StepBytesMethod) Inc(size int) (next bool) {
	*m.i += size
	if *m.i >= len(*m.buf) {
		*m.b = true
		return false
	}
	return true
}

// End breaks the loop
//
// this method always returns false
func (m *StepBytesMethod) End() (next bool) {
	*m.b = true
	return false
}

// Loop creates an inner loop that continues to verify the array length
//
// if loop is not incramented after the callback, the loop will automatically break and log a warning
//
// @cb:
//   - return true to continue loop
//   - return false to break the loop
func (m *StepBytesMethod) Loop(logic func() bool, cb func() bool) {
	lastInd := *m.i
	for *m.i < len(*m.buf) && logic() {
		if !cb() {
			break
		}

		// prevent accidental infinite loop
		if *m.i == lastInd {
			fmt.Println("Warning: Loop Not Incramemted!")
			break
		}
		lastInd = *m.i
	}
}

// Replace will replace bytes at a specific start and end point
//
// this method will also modify the size of the []byte if needed,
// and automatically correct `*i` to the correct index if it changes
func (m *StepBytesMethod) Replace(ind *[2]int, rep *[]byte) {
	*m.buf = append((*m.buf)[:(*ind)[0]], append(*rep, (*m.buf)[(*ind)[1]:]...)...)
	*m.i += ((*ind)[0] - (*ind)[1]) + len(*rep)
}

// GetBuf will return a []byte from the current `*i` index, to the size specified
func (m *StepBytesMethod) GetBuf(size int) []byte {
	if *m.i+size < len(*m.buf) {
		return (*m.buf)[*m.i : *m.i+size]
	}
	return []byte{}
}
