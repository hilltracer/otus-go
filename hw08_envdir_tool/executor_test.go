package main

import "testing"

func TestRunCmd(t *testing.T) {
	env := Environment{
		"ONE":   {Value: "1"},
		"TWO":   {NeedRemove: true},
		"THREE": {Value: ""},
	}

	code := RunCmd([]string{"/bin/bash", "-c", `[ "$ONE" = "1" ] && [ -z "$TWO" ] && [ -z "$THREE" ]`}, env)
	if code != 0 {
		t.Fatalf("command returned %d, want 0", code)
	}
}
