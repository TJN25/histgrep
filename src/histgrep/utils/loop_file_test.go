package utils

import (
	"os"
	"testing"

	"github.com/TJN25/histgrep/hsdata"
)

func TestMain(m *testing.M) {
	InitializeLogger(0)
	os.Exit(m.Run())
}

func TestStartsWithWord(t *testing.T) {
	cases := []struct {
		s, word string
		want    bool
	}{
		{"ls", "ls", true},
		{"ls -la", "ls", true},
		{"ls\t-la", "ls", true},
		{"lsof -i :8080", "ls", false},
		{"cat_tool run", "cat", false},
		{"cat file.txt", "cat", true},
		{"echo", "cd", false},
		{"", "ls", false},
	}
	for _, c := range cases {
		if got := startsWithWord(c.s, c.word); got != c.want {
			t.Errorf("startsWithWord(%q, %q) = %v, want %v", c.s, c.word, got, c.want)
		}
	}
}

func TestFormatLineStartsWithExclude(t *testing.T) {
	format := hsdata.FormattingData{
		Output: map[string][]string{"keys": {"command"}, "separators": {}},
		Excludes: map[string]map[string][]string{
			"command": {"starts_with": {"ls", "cat"}},
		},
	}
	for cmd, kept := range map[string]bool{
		"ls -la":       false,
		"cat foo":      false,
		"lsof -i":      true,
		"cat_tool run": true,
	} {
		terms := MapFormat{"command": cmd}
		got := FormatLine(&terms, &format, true) != ""
		if got != kept {
			t.Errorf("FormatLine(%q) kept = %v, want %v", cmd, got, kept)
		}
	}
}
