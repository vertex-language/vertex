// driver/run.go
package driver

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// execute writes bin to a temp file, runs it with args, streams its output
// straight through, and returns its exit code.
//
// Streaming rather than buffering is the whole reason this exists instead
// of a call to vvm.RunModule: a program that prints as it goes should
// print as it goes, and one that reads argv should get one.
func execute(opts *Options, bin []byte, args []string) (int, error) {
	path, cleanup, err := writeTempExecutable(bin)
	if err != nil {
		return 1, err
	}
	defer cleanup()

	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr

	err = cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), nil
	}
	if err != nil {
		return 1, fmt.Errorf("running the compiled program: %w", err)
	}
	return 0, nil
}

// captureOutput is execute's buffered sibling, used by the test runner:
// a test's whole result is its rendered output, so it has to be compared,
// not streamed.
func captureOutput(bin []byte) (stdout, stderr []byte, code int, err error) {
	path, cleanup, err := writeTempExecutable(bin)
	if err != nil {
		return nil, nil, 1, err
	}
	defer cleanup()

	cmd := exec.Command(path)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	runErr := cmd.Run()
	if exitErr, ok := runErr.(*exec.ExitError); ok {
		return out.Bytes(), errb.Bytes(), exitErr.ExitCode(), nil
	}
	if runErr != nil {
		return out.Bytes(), errb.Bytes(), 1, runErr
	}
	return out.Bytes(), errb.Bytes(), 0, nil
}

// writeTempExecutable is the piece both paths share. The ".exe" suffix on
// Windows is load-bearing: CreateProcess identifies an executable by
// extension and refuses to launch a suffix-less file even when its bytes
// are a valid PE image. ELF and Mach-O need only the +x bit.
func writeTempExecutable(bin []byte) (path string, cleanup func(), err error) {
	pattern := "vertex-run-*"
	if runtime.GOOS == "windows" {
		pattern = "vertex-run-*.exe"
	}

	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", func() {}, fmt.Errorf("creating a temp binary: %w", err)
	}
	path = f.Name()
	cleanup = func() { os.Remove(path) }

	if _, err := f.Write(bin); err != nil {
		f.Close()
		cleanup()
		return "", func() {}, fmt.Errorf("writing the temp binary: %w", err)
	}
	f.Close()
	if err := os.Chmod(path, 0o755); err != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("chmod on the temp binary: %w", err)
	}
	return path, cleanup, nil
}