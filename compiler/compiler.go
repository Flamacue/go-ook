package compiler

import (
	"bytes"
	"fmt"
)

type InsType string

const (
	ADD        InsType = "Ook. Ook."
	SUB        InsType = "Ook! Ook!"
	RIGHT      InsType = "Ook. Ook?"
	LEFT       InsType = "Ook? Ook."
	PUT        InsType = "Ook! Ook."
	READ       InsType = "Ook. Ook!"
	LOOP_BEGIN InsType = "Ook! Ook?"
	LOOP_END   InsType = "Ook? Ook!"
	EOF        InsType = "EOF"
	OOK        string  = "Ook"
)

var ookMap = map[byte](map[byte]InsType){
	'.': map[byte]InsType{
		'.': ADD,
		'!': READ,
		'?': RIGHT,
	},
	'!': map[byte]InsType{
		'!': SUB,
		'.': PUT,
		'?': LOOP_BEGIN,
	},
	'?': map[byte]InsType{
		'.': LEFT,
		'!': LOOP_END,
	},
}

type Instruction struct {
	Type InsType
	Arg  int
}

type Compiler struct {
	code     string
	length   int
	position int

	instructions []*Instruction
}

func New(code string) *Compiler {
	return &Compiler{
		code:         code,
		length:       len(code),
		instructions: []*Instruction{},
	}
}

func (c *Compiler) Compile() ([]*Instruction, error) {
	loopStack := []int{}
	for c.position < c.length {
		// Skip only 1 space inbetween if it exists
		if c.code[c.position] == ' ' {
			c.position++
		}

		// Skip newlines
		for c.position < c.length && c.code[c.position] == '\n' {
			c.position++
		}
		insType, err := c.parseOok(false)
		if err != nil {
			return nil, err
		}

		if insType == EOF {
			break
		}

		switch insType {
		case ADD:
			c.compileFoldable(ADD)
		case SUB:
			c.compileFoldable(SUB)
		case RIGHT:
			c.compileFoldable(RIGHT)
		case LEFT:
			c.compileFoldable(LEFT)
		case PUT:
			c.compileFoldable(PUT)
		case READ:
			c.compileFoldable(READ)
		case LOOP_BEGIN:
			pos := c.newInstruction(LOOP_BEGIN, 0)
			loopStack = append(loopStack, pos)
		case LOOP_END:
			beginLoopIdx := loopStack[len(loopStack)-1]
			loopStack = loopStack[:len(loopStack)-1]

			endLoopIdx := c.newInstruction(LOOP_END, beginLoopIdx)
			c.instructions[beginLoopIdx].Arg = endLoopIdx
		}

		c.position++
	}
	return c.instructions, nil
}

func (c *Compiler) compileFoldable(insType InsType) {
	count := 1
	for nxtInsType, err := c.parseOok(true); err != nil && nxtInsType == insType; nxtInsType, err = c.parseOok(true) {
		c.position += 8
		count++
	}
	c.newInstruction(insType, count)
}

func (c *Compiler) newInstruction(insType InsType, arg int) int {
	ins := &Instruction{Type: insType, Arg: arg}
	c.instructions = append(c.instructions, ins)
	return len(c.instructions) - 1
}

// This function parses the next 2 Ook isntances
func (c *Compiler) parseOok(peek bool) (InsType, error) {
	// Check if we have enough room to read an Ook
	if c.position == c.length {
		return EOF, nil
	}
	var ookIns InsType
	peekPosition := c.position
	if peekPosition+8 > c.length {
		return ookIns, fmt.Errorf("Read position of compiler cannot read the next set of Ooks at %d", c.position)
	}
	var ookBuffer bytes.Buffer
	var first, second byte
	for i := peekPosition + 9; peekPosition < i; peekPosition++ {
		if i-peekPosition == 6 || i-peekPosition == 1 { // We've built our Ook. Check if it's correct and get its type
			if ookBuffer.String() != OOK {
				return ookIns, fmt.Errorf("Non Ook identified at position %d", peekPosition)
			}

			ookBuffer.Reset()

			switch c.code[peekPosition] {
			case '.', '!', '?':
				if first != 0 {
					second = c.code[peekPosition]
					continue
				} else {
					first = c.code[peekPosition]
					continue
				}
			case ' ':
				continue
			default:
				return ookIns, fmt.Errorf("Ook contained bad character at position %d", peekPosition)
			}

		}

		if i-peekPosition == 5 {
			if c.code[peekPosition] != ' ' {
				return ookIns, fmt.Errorf("Expected Whitespace inbetween successive Ooks")
			}
			continue
		}

		ookBuffer.WriteByte(c.code[peekPosition])
	}
	if !peek {
		c.position += 8
	}

	if val, ok := ookMap[first][second]; ok {
		ookIns = val
		return ookIns, nil
	}
	return ookIns, fmt.Errorf("Bad Ook Combination at position %d", peekPosition)

}
