package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(entries))

	for _, e := range entries {
		if e.IsDir() {
			continue // ignore sub-directories
		}

		name := e.Name()
		if strings.ContainsRune(name, '=') {
			return nil, errors.New("environment variable name contains '='")
		}

		path := dir + "/" + name
		info, err := e.Info()
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			env[name] = EnvValue{NeedRemove: true}
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		sc := bufio.NewScanner(f)
		sc.Scan()
		line := sc.Text()
		_ = f.Close()

		// Trim trailing spaces / tabs
		line = strings.TrimRight(line, " \t")
		// Replace 0x00 → '\n'
		line = strings.ReplaceAll(line, "\x00", "\n")

		env[name] = EnvValue{Value: line}
	}
	return env, nil
}
