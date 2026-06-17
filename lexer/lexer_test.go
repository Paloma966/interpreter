package lexer

import (
	"testing"
)

func TestLexSingleInstructions(t *testing.T) {
	tests := []struct {
		input    string
		expected []Instruction
	}{
		{
			"哈",
			[]Instruction{{OpCode: 1, Param: 0, Indent: 0, Line: 1}},
		},
		{
			"哈哈",
			[]Instruction{{OpCode: 2, Param: 0, Indent: 0, Line: 1}},
		},
		{
			"哈哈哈",
			[]Instruction{{OpCode: 3, Param: 0, Indent: 0, Line: 1}},
		},
		{
			"哈哈哈哈",
			[]Instruction{{OpCode: 4, Param: 0, Indent: 0, Line: 1}},
		},
		{
			"哈哈哈哈哈",
			[]Instruction{{OpCode: 5, Param: 0, Indent: 0, Line: 1}},
		},
		{
			"哈哈哈哈哈哈",
			[]Instruction{{OpCode: 6, Param: 0, Indent: 0, Line: 1}},
		},
		{
			"哈哈哈哈哈哈哈",
			[]Instruction{{OpCode: 7, Param: 0, Indent: 0, Line: 1}},
		},
		{
			"哈哈哈哈哈哈哈哈",
			[]Instruction{{OpCode: 8, Param: 0, Indent: 0, Line: 1}},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)
		instructions := l.Lex()

		if len(l.Errors()) > 0 {
			t.Errorf("input %q: unexpected errors: %v", tt.input, l.Errors())
			continue
		}

		if len(instructions) != len(tt.expected) {
			t.Errorf("input %q: expected %d instructions, got %d", tt.input, len(tt.expected), len(instructions))
			continue
		}

		for i, inst := range instructions {
			if inst != tt.expected[i] {
				t.Errorf("input %q[%d]: expected %+v, got %+v", tt.input, i, tt.expected[i], inst)
			}
		}
	}
}

func TestLexWithParams(t *testing.T) {
	input := "哈   \n哈哈 \n哈哈哈  "
	l := New(input)
	instructions := l.Lex()

	if len(l.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", l.Errors())
	}

	if len(instructions) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(instructions))
	}

	if instructions[0].Param != 3 {
		t.Errorf("expected param 3, got %d", instructions[0].Param)
	}
	if instructions[1].Param != 1 {
		t.Errorf("expected param 1, got %d", instructions[1].Param)
	}
	if instructions[2].Param != 2 {
		t.Errorf("expected param 2, got %d", instructions[2].Param)
	}
}

func TestLexWithIndent(t *testing.T) {
	input := "哈\n  哈哈\n    哈哈哈"
	l := New(input)
	instructions := l.Lex()

	if len(l.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", l.Errors())
	}

	if len(instructions) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(instructions))
	}

	if instructions[0].Indent != 0 {
		t.Errorf("expected indent 0, got %d", instructions[0].Indent)
	}
	if instructions[1].Indent != 2 {
		t.Errorf("expected indent 2, got %d", instructions[1].Indent)
	}
	if instructions[2].Indent != 4 {
		t.Errorf("expected indent 4, got %d", instructions[2].Indent)
	}
}

func TestLexSkipsEmptyLines(t *testing.T) {
	input := "哈\n\n哈哈\n\n哈哈哈"
	l := New(input)
	instructions := l.Lex()

	if len(l.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", l.Errors())
	}

	if len(instructions) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(instructions))
	}
}

func TestLexTooManyHa(t *testing.T) {
	input := "哈哈哈哈哈哈哈哈哈" // 9 哈
	l := New(input)
	l.Lex()

	if len(l.Errors()) == 0 {
		t.Fatal("expected error for too many 哈, got none")
	}
}

func TestLexIllegalChar(t *testing.T) {
	input := "哈a"
	l := New(input)
	l.Lex()

	if len(l.Errors()) == 0 {
		t.Fatal("expected error for illegal char, got none")
	}
}

func TestLexSpacesOnlyLine(t *testing.T) {
	input := "   \n  \n哈" // spaces-only lines are skipped like empty lines
	l := New(input)
	instructions := l.Lex()

	if len(l.Errors()) > 0 {
		t.Fatalf("unexpected errors: %v", l.Errors())
	}
	if len(instructions) != 1 {
		t.Fatalf("expected 1 instruction, got %d", len(instructions))
	}
}
