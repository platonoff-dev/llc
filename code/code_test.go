package code

import (
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("[%d]", i)
		t.Run(name, func(t *testing.T) {
			instruction := Make(tt.op, tt.operands...)
			if len(instruction) != len(tt.expected) {
				t.Errorf("instruction has wrong length. got=%d, expected=%d", len(instruction), len(tt.expected))
			}

			for i, b := range tt.expected {
				if instruction[i] != tt.expected[i] {
					fmt.Errorf("wrong byte at pos %d got=%d, expected=%d", i, instruction[i], b)
				}
			}
		})
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected {
		t.Errorf("instructions wrongly formated.\nwant=%q\ngot=%q", expected, concatted.String())
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("[%d]", i)
		t.Run(name, func(t *testing.T) {
			instruction := Make(tt.op, tt.operands...)

			def, err := Lookup(byte(tt.op))
			if err != nil {
				t.Fatalf("definition not found: %q\n", err)
			}

			operandsRead, n := ReadOperands(def, instruction[1:])
			if n != tt.bytesRead {
				t.Fatalf("n wrong. want=%d, got=%d", tt.bytesRead, n)
			}

			for i, want := range tt.operands {
				if operandsRead[i] != want {
					t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
				}
			}
		})
	}
}
