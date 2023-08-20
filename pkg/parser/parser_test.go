package parser

import "testing"

func TestCommand1(t *testing.T) {
	parser := NewParser("say -m=world")
	args := parser.Parse()

	value, ok := args.GetFlag("m")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "world" {
		t.Fatal("Flag value set incorretly.")
	}
}
