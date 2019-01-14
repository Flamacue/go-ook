package compiler

import (
	"testing"
)

//TODO: Make much better tests

func TestParseOok(t *testing.T) {
	input := `Ook. Ook!`

	c := New(input)

	insType, err := c.parseOok(false)
	if err != nil {
		t.Fatalf("err was not nil. %q", err)
	}
	if insType != READ {
		t.Fatalf("Instruction type incorrect. got=%q. expected=READ", insType)
	}

}

func TestCompile(t *testing.T) {
	input := `Ook. Ook! Ook. Ook!`

	c := New(input)

	instructions, err := c.Compile()
	if err != nil {
		t.Fatalf("err was not nil. %q", err)
	}
	if len(instructions) == 1 {
		t.Fatalf("Did not make correct number of instructions")
	}
}
