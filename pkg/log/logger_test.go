package log

import "testing"

func Test_padLevelName(t *testing.T) {

	expectations := map[string]string{
		"TRACE": "TRACE",
		"DEBUG": "DEBUG",
		"INFO": "INFO ",
		"WARN": "WARN ",
		"ERROR": "ERROR",
		"FATAL": "FATAL",
		"FOO": "FOO  ",
		"G": "G    ",
	}

	for input, want := range expectations {
		got := padLevelName(input)

		if got != want {
			t.Errorf("Input: %s, Got: %s, Want: %s", input, want, got)
		}
	}
}
