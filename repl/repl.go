package repl

import (
	"bufio"
	"fmt"
	"interpreter/evaluator"
	"interpreter/lexer"
	"io"
	"strings"
)

const PROMPT = "哈> "
const CONT = "…> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	var lines []string

	for {
		if len(lines) == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT)
		}

		scanned := scanner.Scan()
		if !scanned {
			// EOF: execute any remaining program
			if len(lines) > 0 {
				runProgram(strings.Join(lines, "\n"), out)
			}
			return
		}

		line := scanner.Text()

		// Empty line means execute accumulated program
		if strings.TrimSpace(line) == "" {
			if len(lines) > 0 {
				runProgram(strings.Join(lines, "\n"), out)
				lines = nil
			}
			continue
		}

		lines = append(lines, line)
	}
}

func runProgram(input string, out io.Writer) {
	l := lexer.New(input)
	program := l.Lex()

	if len(l.Errors()) > 0 {
		for _, e := range l.Errors() {
			fmt.Fprintln(out, e)
		}
		return
	}

	if len(program) == 0 {
		return
	}

	m := evaluator.New(program)
	if err := m.Run(); err != nil {
		fmt.Fprintln(out, "错误:", err)
		return
	}

	output := m.Output()
	if len(output) > 0 {
		fmt.Fprintln(out, output)
	}
}
