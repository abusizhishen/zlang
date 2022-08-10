package parser

import (
	"fmt"
	"github.com/abusizhishen/zlang/ast"
	"github.com/abusizhishen/zlang/lexer"
	"testing"
)

func TestParser_ParseStatement(t *testing.T) {
	//testLetStatements(t)
	testOperatorPrecedenceParsing(t)
}

func testLetStatements(t *testing.T) {
	input := `
	let a= 3;
let b=5;
let  foobar=8000;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatal("parseProgram err,nil")
	}

	checkErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("p.statements does not contain 3 statement, got:%d", len(program.Statements))
	}

	tests := []struct {
		ExpectIdentifier string
	}{
		{"a"},
		{"b"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.ExpectIdentifier) {
			return
		}
	}

	s := program.String()
	fmt.Println(s)
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLIteral not let,got:%q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got: %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.name.value not %s.got:%q", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.name.TokenLiteral not %s. got:%q", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parse has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parse error: %q", err)
	}

	t.FailNow()
}

func testOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"let a = if (x){x}else{y}", ""},
		{"(2+(3+4))+1", "((2 + (3 + 4)) + 1)"},
		{"(2+3)+1", "((2 + 3) + 1)"},
		{"1+(2+3)", "(1 + (2 + 3))"},
		{"1+(2+3)+4", "((1 + (2 + 3)) + 4)"},
		{"true == false", "(true == false)"},
		{"a+b+c", "((a + b) + c)"},
		{"a+b*c", "(a + (b * c))"},
		{"a*b*c", "((a * b) * c)"},
		{"a+b/c", "(a + (b / c))"},
		{"a+b*c+d/e-f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3+4;-5*5", "(3 + 4)((-5) * 5)"},
		{"5<4 == 3<4", "((5 < 4) == (3 < 4))"},
		{"5<4 != 3>4", "((5 < 4) != (3 > 4))"},
		{"3+4*5 == 3*1+4*5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3+1-2", "((3 + 1) - 2)"},
		{"1 == 3-2", "(1 == (3 - 2))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected: %q, got=%q input:%q", tt.expected, actual, tt.input)
		} else {
			//t.Log(actual)
		}
	}
}
