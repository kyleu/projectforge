package util

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/buildkite/shellwords"
	"github.com/pkg/errors"
)

func StartProcess(cmd string, path string, in io.Reader, out io.Writer, er io.Writer, env ...string) (*exec.Cmd, error) {
	args, err := shellwords.Split(cmd)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, errors.New("no arguments provided")
	}
	firstArg := args[0]

	if !strings.Contains(firstArg, "/") && !strings.Contains(firstArg, "\\") {
		firstArg, err = exec.LookPath(firstArg)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to look up cmd [%s]", firstArg)
		}
	}

	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	if er == nil {
		er = os.Stderr
	}

	c := &exec.Cmd{Path: firstArg, Args: args, Env: env, Stdin: in, Stdout: out, Stderr: er, Dir: path}
	err = c.Start()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to start [%s] (%T)", cmd, err)
	}
	return c, nil
}

func RunProcess(cmd string, path string, in io.Reader, out io.Writer, er io.Writer, env ...string) (int, error) {
	c, err := StartProcess(cmd, path, in, out, er, env...)
	if err != nil {
		return -1, err
	}
	err = c.Wait()
	if err != nil {
		ec, ok := err.(*exec.ExitError) //nolint:errorlint
		if ok {
			return ec.ExitCode(), nil
		}
		return -1, errors.Wrapf(err, "unable to run [%s] (%T)", cmd, err)
	}
	return 0, nil
}

func RunProcessSimple(cmd string, path string) (int, string, error) {
	var buf bytes.Buffer
	ec, err := RunProcess(cmd, path, nil, &buf, &buf)
	if err != nil {
		return -1, "", err
	}
	return ec, buf.String(), nil
}
