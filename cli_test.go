package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runCliForTesting(t *testing.T, args ...string) (code int, stdout string, stderr string) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	var stdoutBuffer, stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer
	err := cmd.Run()
	if err != nil {
		if _, isExitError := err.(*exec.ExitError); isExitError {
			// must be caused by exit code, ignore the error
		} else {
			t.Fatalf("Failed to run with arguments '%s': %s\n", strings.Join(args, " "), err)
		}
	}
	stdout, stderr = stdoutBuffer.String(), stderrBuffer.String()
	code = cmd.ProcessState.ExitCode()

	stderr = strings.TrimSuffix(stderr, "exit status 1\n")
	return
}

func shouldTestPrintHelp(t *testing.T, stdout string, stderr string) {
	assert.Regexp(t, "^Usage:.*", stdout, "Should print usage to stdout")
	assert.Regexp(t, "Options:", stdout, "Should print options to stdout")
	assert.Regexp(t, "References:", stdout, "Should print references to stdout")
}

func TestCliNoArgument(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t)
	assert.Equal(t, 1, code, "Exit code should be 1")
	shouldTestPrintHelp(t, stdout, stderr)
	assert.Equal(t, "Error: flag file or check file must be specified\n", stderr)
}

func TestCliErrorPanicNonExistingDir(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-f", "package.json", "--error", "panic", "./non-existing-dir")
	assert.Equal(t, 1, code, "Exit code should be 1")
	assert.Equal(t, 0, len(stdout), "Should have nothing sent to stdout")
	assert.Equal(t, "Error: stat ./non-existing-dir: no such file or directory\n", stderr)
}

func TestCliErrorIgnoreNonExistingDir(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-f", "package.json", "--error", "ignore", "./non-existing-dir")
	assert.Equal(t, 0, code, "Exit code should be 0")
	assert.Equal(t, 0, len(stdout), "Should have nothing sent to stdout")
	assert.Equal(t, 0, len(stderr), "Should have nothing sent to stderr")
}

func TestCliErrorPrintNonExistingDir(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-f", "package.json", "--error", "print", "./non-existing-dir")
	assert.Equal(t, 0, code, "Exit code should be 0")
	assert.Equal(t, 0, len(stdout), "Should have nothing sent to stdout")
	assert.Equal(t, "Error: stat ./non-existing-dir: no such file or directory\n", stderr)
}

func TestCliHelp(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-h")
	assert.Equal(t, 0, code, "Exit code should be 0")
	shouldTestPrintHelp(t, stdout, stderr)
	assert.Equal(t, "", stderr)
}

func TestCliPackageJsonInRepo1(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-f", "package.json", "testdata/repo1")
	assert.Equal(t, 0, code, "Exit code should be 0")
	assert.Equal(t, "", stderr)
	assert.Equal(t,
		`testdata/repo1
testdata/repo1/inbound
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/china
testdata/repo1/outbound/china/mainland
`,
		stdout, "Should output exactly these")
}

func TestCliPrint0PackageJsonInRepo1(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-0", "-f", "package.json", "testdata/repo1")
	assert.Equal(t, 0, code, "Exit code should be 0")
	assert.Equal(t, "", stderr)
	assert.Equal(t,
		strings.ReplaceAll(`testdata/repo1
testdata/repo1/inbound
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/china
testdata/repo1/outbound/china/mainland
`, "\n", string([]byte{0})),
		stdout, "Should output exactly these")
}

func TestCliDepth0PackageJsonInRepo1(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-f", "package.json", "-d", "0", "testdata/repo1")
	assert.Equal(t, 0, code, "Exit code should be 0")
	assert.Equal(t, "", stderr)
	assert.Equal(t, "testdata/repo1\n", stdout, "Should output the root directory")
}

func TestCliNoDefaultExcludesPackageJsonInRepo1(t *testing.T) {
	code, stdout, stderr := runCliForTesting(t, "-f", "package.json", "--no-default-excludes", "testdata/repo1")
	assert.Equal(t, 0, code, "Exit code should be 0")
	assert.Equal(t, "", stderr)
	assert.Equal(t,
		`testdata/repo1
testdata/repo1/inbound
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/china
testdata/repo1/outbound/china/mainland
testdata/repo1/outbound/china/mainland/node_modules/package1
testdata/repo1/outbound/china/mainland/node_modules/package2
`,
		stdout, "Should output exactly these")
}
