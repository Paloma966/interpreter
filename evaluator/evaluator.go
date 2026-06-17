package evaluator

import (
	"bytes"
	"fmt"
	"interpreter/lexer"
	"io"
	"os"
)

type Machine struct {
	tape    [30000]uint8
	ptr     int
	program []lexer.Instruction
	ip      int
	output  bytes.Buffer
	input   io.Reader
	jumps   map[int]int // ip -> target ip for [ and ]
}

func New(program []lexer.Instruction) *Machine {
	m := &Machine{
		program: program,
		input:   os.Stdin,
		jumps:   buildJumps(program),
	}
	return m
}

func (m *Machine) SetInput(r io.Reader) {
	m.input = r
}

func (m *Machine) Output() string {
	return m.output.String()
}

func (m *Machine) Run() error {
	for m.ip < len(m.program) {
		inst := m.program[m.ip]

		switch inst.OpCode {
		case 1: // >
			n := inst.Param
			if n == 0 {
				n = 1
			}
			m.ptr += n
			if m.ptr >= len(m.tape) {
				return fmt.Errorf("行%d: 指针越界(右)", inst.Line)
			}

		case 2: // <
			n := inst.Param
			if n == 0 {
				n = 1
			}
			m.ptr -= n
			if m.ptr < 0 {
				return fmt.Errorf("行%d: 指针越界(左)", inst.Line)
			}

		case 3: // +
			n := inst.Param
			if n == 0 {
				n = 1
			}
			m.tape[m.ptr] += uint8(n)

		case 4: // -
			n := inst.Param
			if n == 0 {
				n = 1
			}
			m.tape[m.ptr] -= uint8(n)

		case 5: // .
			m.output.WriteByte(m.tape[m.ptr])

		case 6: // ,
			var buf [1]byte
			_, err := m.input.Read(buf[:])
			if err != nil {
				m.tape[m.ptr] = 0
			} else {
				m.tape[m.ptr] = buf[0]
			}

		case 7: // [
			if m.tape[m.ptr] == 0 {
				m.ip = m.jumps[m.ip]
			}

		case 8: // ]
			if m.tape[m.ptr] != 0 {
				m.ip = m.jumps[m.ip]
			}
		}

		m.ip++
	}

	return nil
}

func buildJumps(program []lexer.Instruction) map[int]int {
	jumps := make(map[int]int)

	// For each indentation level, find matching 7/8 pairs
	// Use a stack: push opcode 7 ips, pop on matching opcode 8
	type frame struct {
		ip     int
		indent int
	}
	var stack []frame

	for i, inst := range program {
		if inst.OpCode == 7 {
			stack = append(stack, frame{ip: i, indent: inst.Indent})
		} else if inst.OpCode == 8 {
			if len(stack) == 0 {
				continue // error: unmatched ], will cause runtime issue
			}
			top := stack[len(stack)-1]
			if top.indent == inst.Indent {
				// Match found: jump to the bracket itself (ip++ handles the advance)
				jumps[top.ip] = i // [ → ]
				jumps[i] = top.ip // ] → [
				stack = stack[:len(stack)-1]
			} else {
				// Indentation mismatch — unwind stack until match
				for len(stack) > 0 && stack[len(stack)-1].indent != inst.Indent {
					stack = stack[:len(stack)-1]
				}
				if len(stack) > 0 {
					top = stack[len(stack)-1]
					jumps[top.ip] = i
					jumps[i] = top.ip
					stack = stack[:len(stack)-1]
				}
			}
		}
	}

	return jumps
}
