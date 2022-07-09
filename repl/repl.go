package repl

import (
	"io"
	"zlang/lexer"
	"zlang/token"
)
import "bufio"
import "fmt"

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scaner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		scan := scaner.Scan()
		if !scan {
			return
		}

		line := scaner.Text()
		lex := lexer.New(line)
		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
