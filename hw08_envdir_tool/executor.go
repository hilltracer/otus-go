package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1 // nothing to run
	}

	// Build resulting environment as map for easy overwrite/removal
	base := map[string]string{}
	for _, v := range os.Environ() {
		parts := strings.SplitN(v, "=", 2)
		base[parts[0]] = parts[1]
	}

	for k, v := range env {
		if v.NeedRemove {
			delete(base, k)
		} else {
			base[k] = v.Value
		}
	}

	// Convert back to []string
	finalEnv := make([]string, 0, len(base))
	for k, v := range base {
		finalEnv = append(finalEnv, k+"="+v)
	}
	// #nosec G204 -- we are proxy, all responsibility is on the user
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = finalEnv
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
