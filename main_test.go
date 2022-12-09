package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupFlags()
}

// Run doMain(...)
//
// Results:
//   - output: the content supposed to be sent to stdout
//   - error: the content supposed to be sent to stderr
func runDoMainForTesting(args ...string) (output string, error string) {
	// setup arguments
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = append([]string{"cmd"}, args...)

	// clear up output
	output = ""
	error = ""

	// doMain()
	doMain(func(text string) {
		output += text
	}, func(errors []string, printUsage bool, exitCode int) {
		if len(errors) > 0 {
			for _, errorString := range errors {
				error += "Error: " + errorString + "\n"
			}
		} else {
			error = ""
		}
	})
	return
}

func TestDoMainNoArgument(t *testing.T) {
	output, error := runDoMainForTesting()
	assert.Equal(t, "", output)
	assert.Greater(t, len(error), 0, "There should be error")
	assert.Equal(t, "Error: flag file has not been specified\n", error)
}

func TestDoMainHelp(t *testing.T) {
	output, error := runDoMainForTesting("--help")
	assert.Equal(t, 0, len(error), "There should be no error")
	assert.Equal(t, "", output)
}

var validArgumentsAndExpectedOutputs = []struct {
	arguments string
	output    string
}{
	{
		"-f package.json testdata/repo1",
		`testdata/repo1
testdata/repo1/inbound
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/china
testdata/repo1/outbound/china/mainland
`,
	}, {
		"-f package.json -d 0 testdata/repo1",
		"testdata/repo1\n",
	}, {
		"-f package.json --no-default-excludes testdata/repo1",
		`testdata/repo1
testdata/repo1/inbound
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/china
testdata/repo1/outbound/china/mainland
testdata/repo1/outbound/china/mainland/node_modules/package1
testdata/repo1/outbound/china/mainland/node_modules/package2
`,
	}, {
		"-f package.json --depth 0 testdata/repo1/inbound",
		"testdata/repo1/inbound\n",
	}, {
		"-f package.json --depth 0 --subdirectories-only --no-default-excludes testdata/repo1/inbound",
		"",
	}, {
		"-f package.json --subdirectories-only --no-default-excludes testdata/repo1/outbound/china/mainland",
		`testdata/repo1/outbound/china/mainland/node_modules/package1
testdata/repo1/outbound/china/mainland/node_modules/package1/node_modules/package1-1
testdata/repo1/outbound/china/mainland/node_modules/package2
`,
	}, {
		"-f serverless.* testdata/repo1",
		`testdata/repo1/inbound
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/australia
testdata/repo1/outbound/china/sars
`,
	}, {
		"-f serverless.ts testdata/repo1",
		"testdata/repo1/outbound/china/sars\n",
	}, {
		"-f serverless.* -c build.gradle testdata/repo1",
		`testdata/repo1/outbound/australia
testdata/repo1/outbound/china/sars
`,
	}, {
		"-f package.* -c package.json testdata/repo1",
		`testdata/repo1
testdata/repo1/inbound
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/china
testdata/repo1/outbound/china/mainland
`,
	}, {
		"-f package.* -c package.yml testdata/repo1",
		`testdata/repo1/api
`,
	}, {
		`-f package.* -c package.json -e "@types/mocha": testdata/repo1`,
		`testdata/repo1/inbound
`,
	}, {
		`-f package.* -c package.json -e "@types/mocha": -i testdata/repo1`,
		`testdata/repo1
testdata/repo1/outbound/New Zealand
testdata/repo1/outbound/china
testdata/repo1/outbound/china/mainland
`,
	}, {
		`-f package.* -c package.json -e "volta": testdata/repo1`,
		`testdata/repo1/inbound
`,
	}, {
		`-f package.* -c package.json -e "dependencies":\s*{[^{}]*"volta": testdata/repo1`,
		"",
	}, {
		`-f package.* -c package.json -e "dependencies":\s*{[^{}]*"mocha": testdata/repo1`,
		`testdata/repo1/inbound
`,
	}, {
		`-f build.gradle* testdata/repo1`,
		`testdata/repo1/outbound/australia
testdata/repo1/outbound/china/sars
testdata/repo1/outbound/usa
`,
	}, {
		`-f build.gradle testdata/repo1`,
		`testdata/repo1/outbound/australia
testdata/repo1/outbound/china/sars
`,
	}, {
		`-f build.gradle -f mvn.* testdata/repo1`,
		`testdata/repo1/outbound/australia
testdata/repo1/outbound/china/sars
testdata/repo1/storage
`,
	}, {
		`-f build.gradle -f mvn.* -x **/australia -x **/storage testdata/repo1`,
		`testdata/repo1/outbound/china/sars
`,
	}, {
		`-f anything testdata/non-existing-dir`,
		"",
	}, {
		`-f anything --error ignore testdata/non-existing-dir`,
		"",
	},
}

var invalidArgumentsAndExpectedOutputs = []struct {
	arguments string
	output    string
	error     string
}{
	{
		`testdata/non-existing-dir`,
		"",
		"Error: flag file has not been specified\n",
	},
	{
		`-f anything --error print testdata/non-existing-dir`,
		"",
		"Error: stat testdata/non-existing-dir: no such file or directory\n",
	},
	{
		`-f anything --error panic testdata/non-existing-dir`,
		"",
		"Error: stat testdata/non-existing-dir: no such file or directory\n",
	},
}

func TestDoMainWithValidArguments(t *testing.T) {
	spaces := regexp.MustCompile(" +")
	for _, vaaeo := range validArgumentsAndExpectedOutputs {
		args := spaces.Split(vaaeo.arguments, -1)
		t.Run(vaaeo.arguments, func(t *testing.T) {
			output, error := runDoMainForTesting(args...)
			assert.Equal(t, 0, len(error), "There should be no error")
			assert.Equal(t, vaaeo.output, output, "Should output exactly these")
		})
	}
}

func TestDoMainWithInvalidArguments(t *testing.T) {
	spaces := regexp.MustCompile(" +")
	for _, vaaeo := range invalidArgumentsAndExpectedOutputs {
		args := spaces.Split(vaaeo.arguments, -1)
		t.Run(vaaeo.arguments, func(t *testing.T) {
			output, error := runDoMainForTesting(args...)
			assert.Equal(t, vaaeo.error, error, "Should generate exactly these error")
			assert.Equal(t, vaaeo.output, output, "Should output exactly these")
		})
	}
}
