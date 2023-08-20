package parser

import "testing"

func TestSettingShortOption(t *testing.T) {
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

func TestSettingManyShortOption(t *testing.T) {
	parser := NewParser(`say -m world -c "black" -v`)
	args := parser.Parse()

	value, ok := args.GetFlag("m")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "world" {
		t.Fatal("Flag value set incorretly.")
	}

	value, ok = args.GetFlag("c")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "black" {
		t.Fatal("Flag value set incorretly.")
	}

	_, ok = args.GetFlag("v")

	if !ok {
		t.Fatal("Flag not set properly!")
	}
}

func TestSettingLongOption(t *testing.T) {
	parser := NewParser(`say --message=world --colour "white" --verbose false --state="none"`)
	args := parser.Parse()

	value, ok := args.GetFlag("message")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "world" {
		t.Fatal("Flag value set incorretly.")
	}

	value, ok = args.GetFlag("colour")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "white" {
		t.Fatal("Flag value set incorretly.")
	}

	value, ok = args.GetFlag("verbose")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "false" {
		t.Fatal("Flag value set incorretly.")
	}

	value, ok = args.GetFlag("state")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "none" {
		t.Fatal("Flag value set incorretly.")
	}
}

func TestSettingMultipleLongOption(t *testing.T) {
	parser := NewParser("say --FLAG1 --FLAG2 --FLAG3")
	args := parser.Parse()

	_, ok := args.GetFlag("FLAG1")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	_, ok = args.GetFlag("FLAG2")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	_, ok = args.GetFlag("FLAG3")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

}

func TestSettingManyLongOption(t *testing.T) {
	parser := NewParser("say --message=world --colour black --verbose")
	args := parser.Parse()

	value, ok := args.GetFlag("message")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "world" {
		t.Fatal("Flag value set incorretly.")
	}

	value, ok = args.GetFlag("colour")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "black" {
		t.Fatal("Flag value set incorretly.")
	}

	_, ok = args.GetFlag("verbose")

	if !ok {
		t.Fatal("Flag not set properly!")
	}
}

func TestSettingManyLongOption2(t *testing.T) {
	parser := NewParser(`say --message="Hello world!" --colour="black" --verbose`)
	args := parser.Parse()

	value, ok := args.GetFlag("message")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "Hello world!" {
		t.Fatal("Flag value set incorretly.")
	}

	value, ok = args.GetFlag("colour")

	if !ok {
		t.Fatal("Flag not set properly!")
	}

	if value != "black" {
		t.Fatal("Flag value set incorretly.")
	}

	_, ok = args.GetFlag("verbose")

	if !ok {
		t.Fatal("Flag not set properly!")
	}
}
