package parser

import "testing"

func TestIsWhiteSpace(t *testing.T) {
	if !isWhiteSpace(rune(' ')) {
		t.Fatalf("isWhiteSpace(rune(' ')) = false, want true")
	}

	if !isWhiteSpace(rune('\t')) {
		t.Fatalf("isWhiteSpace(rune('\\t')) = false, want true")
	}

	// We need new line to be false because we are ending commands
	// with a newline (there is no ending character like ;
	if isWhiteSpace(rune('\n')) {
		t.Fatalf("isWhiteSpace(rune('\n')) = true, want false")
	}
}

func TestIsLetter(t *testing.T) {
	if isLetter(rune(' ')) {
		t.Fatalf("isLetter(rune(' ')) = true, want false")
	}

	if isLetter(rune('\t')) {
		t.Fatalf("isLetter(rune('\\t')) = true, want false")
	}

	if isLetter(rune('\n')) {
		t.Fatalf("isLetter(rune('\n')) = true, want false")
	}

	if isLetter(rune('0')) {
		t.Fatalf("isLetter(rune('0')) = true, want false")
	}

	if !isLetter(rune('a')) {
		t.Fatalf("isLetter(rune('a')) = false, want true")
	}

	if !isLetter(rune('A')) {
		t.Fatalf("isLetter(rune('a')) = false, want true")
	}

	if !isLetter(rune('g')) {
		t.Fatalf("isLetter(rune('a')) = false, want true")
	}

	if !isLetter(rune('G')) {
		t.Fatalf("isLetter(rune('a')) = false, want true")
	}
}
