package main

import (
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir := t.TempDir()

	mustWrite := func(name, content string) {
		if err := os.WriteFile(dir+"/"+name, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("FOO", "bar  \t")         // trailing ws should be trimmed
	mustWrite("ZERO", "")               // empty file
	mustWrite("NL", "line1\x00line2\n") // 0x00 → '\n'

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := env["FOO"]; v.Value != "bar" || v.NeedRemove {
		t.Errorf("FOO parsed wrongly: %+v", v)
	}
	if !env["ZERO"].NeedRemove {
		t.Errorf("ZERO should be marked NeedRemove")
	}
	if v := env["NL"].Value; v != "line1\nline2" {
		t.Errorf("NL replacement wrong: %q", v)
	}
}
