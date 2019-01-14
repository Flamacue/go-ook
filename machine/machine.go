package machine

import (
	"io"

	"github.com/flamacue/go-ook/compiler"
)

type Machine struct {
	instructions []*compiler.Instruction
	iptr         int

	tape [30000]int // Since ook is essentially BF
	tptr int

	input  io.Reader
	output io.Writer

	buffer []byte
}

func New(input []*compiler.Instruction, in io.Reader, out io.Writer) *Machine {
	return &Machine{
		instructions: input,
		input:        in,
		output:       out,
		buffer:       make([]byte, 1),
	}
}

func (m *Machine) getChar() {
	num, err := m.input.Read(m.buffer)
	if err != nil {
		panic(err) //TODO: fix this up
	}

	if num != 1 {
		panic("Wrong number of bytes read")
	}

	m.tape[m.tptr] = int(m.buffer[0])
}

func (m *Machine) putChar() {
	m.buffer[0] = byte(m.tape[m.tptr])
	n, err := m.output.Write(m.buffer)
	if err != nil {
		panic(err) //TODO: fix this up
	}

	if n != 1 {
		panic("Wrong number of bytes written")
	}
}

func (m *Machine) Execute() {

	for m.iptr < len(m.instructions) {
		ins := m.instructions[m.iptr]

		switch ins.Type {
		case compiler.ADD:
			m.tape[m.tptr] += ins.Arg
		case compiler.SUB:
			m.tape[m.tptr] -= ins.Arg
		case compiler.RIGHT:
			m.tptr += ins.Arg
		case compiler.LEFT:
			m.tptr -= ins.Arg
		case compiler.PUT:
			for i := 0; i < ins.Arg; i++ {
				m.putChar()
			}
		case compiler.READ:
			for i := 0; i < ins.Arg; i++ {
				m.getChar()
			}
		case compiler.LOOP_BEGIN:
			if m.tape[m.tptr] == 0 {
				m.iptr = ins.Arg
			}
		case compiler.LOOP_END:
			if m.tape[m.tptr] != 0 {
				m.iptr = ins.Arg
			}
		}

		m.iptr++
	}
}
