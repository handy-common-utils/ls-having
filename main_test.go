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

func runDoMainForTesting(args ...string) (output string, hasError bool, error string) {
	// setup arguments
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = append([]string{"cmd"}, args...)

	// clear up output
	output = ""
	hasError = false
	error = ""

	// doMain()
	doMain(func(text string) {
		output += text + "\n"
	}, func(text string) {
		if len(text) > 0 {
			hasError = true
			error = text
		} else {
			hasError = false
			error = ""
		}
	})
	return
}

func TestDoMainNoArgument(t *testing.T) {
	output, hasError, error := runDoMainForTesting()
	assert.Equal(t, "", output)
	assert.True(t, hasError, "There should be error")
	assert.Equal(t, "flag file has not been specified", error)
}

func TestDoMainHelp(t *testing.T) {
	output, hasError, _ := runDoMainForTesting("--help")
	assert.False(t, hasError, "There should be no error")
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
	},
}

func TestDoMainWithValidArguments(t *testing.T) {
	spaces := regexp.MustCompile(" +")
	for _, vaaeo := range validArgumentsAndExpectedOutputs {
		args := spaces.Split(vaaeo.arguments, -1)
		t.Run(vaaeo.arguments, func(t *testing.T) {
			output, hasError, _ := runDoMainForTesting(args...)
			assert.False(t, hasError, "There should be no error")
			assert.Equal(t, vaaeo.output, output, "Should output exactly these")
		})
	}
}
