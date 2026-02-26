package execx

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

func Run(name string, args []string, env []string, timeoutSec int) Result {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	if len(env) > 0 {
		cmd.Env = append(cmd.Env, env...)
	}
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	start := time.Now()
	err := cmd.Run()
	res := Result{Stdout: out.String(), Stderr: errBuf.String(), DurationMS: time.Since(start).Milliseconds()}
	if err == nil {
		res.ExitCode = 0
		return res
	}
	if ee, ok := err.(*exec.ExitError); ok {
		res.ExitCode = ee.ExitCode()
		return res
	}
	res.ExitCode = 124
	return res
}

func RunShell(command string, timeoutSec int) Result {
	return Run("bash", []string{"-lc", command}, nil, timeoutSec)
}
