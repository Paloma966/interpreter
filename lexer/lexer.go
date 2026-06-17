package lexer

import (
	"fmt"
	"strings"
)

type Instruction struct {
	OpCode int // 1-8: number of consecutive 哈
	Param  int // trailing spaces after 哈 (0 = default 1 for ops 1-4)
	Indent int // leading spaces before 哈 (indentation level)
	Line   int // source line number
}

type Lexer struct {
	lines  []string
	errors []string
}

func New(input string) *Lexer {
	return &Lexer{lines: strings.Split(input, "\n")}
}

func (l *Lexer) Lex() []Instruction {
	var program []Instruction

	for lineNum, rawLine := range l.lines {
		// Trim trailing spaces for parsing (but we count leading/trailing separately)
		line := rawLine

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		inst := l.parseLine(line, lineNum+1)
		if inst != nil {
			program = append(program, *inst)
		}
	}

	return program
}

func (l *Lexer) parseLine(line string, lineNum int) *Instruction {
	runes := []rune(line)

	// Phase 1: count leading spaces (indent)
	indent := 0
	pos := 0
	for pos < len(runes) && runes[pos] == ' ' {
		indent++
		pos++
	}

	// Phase 2: count consecutive 哈
	if pos >= len(runes) {
		return nil // shouldn't happen since we skip empty lines
	}

	haCount := 0
	for pos < len(runes) && runes[pos] == '哈' {
		haCount++
		pos++
	}

	if haCount == 0 {
		l.addError(lineNum, "没有哈，啥也不是")
		return nil
	}

	if haCount > 8 {
		l.addError(lineNum, fmt.Sprintf("哈太多了(%d个)，最多8个", haCount))
		return nil
	}

	// Phase 3: count trailing spaces (param)
	param := 0
	for pos < len(runes) && runes[pos] == ' ' {
		param++
		pos++
	}

	// Phase 4: check for illegal characters
	if pos < len(runes) {
		l.addError(lineNum, fmt.Sprintf("非法字符: %q", runes[pos]))
		return nil
	}

	return &Instruction{
		OpCode: haCount,
		Param:  param,
		Indent: indent,
		Line:   lineNum,
	}
}

func (l *Lexer) addError(line int, msg string) {
	l.errors = append(l.errors, fmt.Sprintf("行%d: %s", line, msg))
}

func (l *Lexer) Errors() []string {
	return l.errors
}
