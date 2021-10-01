package main

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// it's possible to test the program this way, however it's not the easiest way to do that,
// so business logic tests are in record package.
func TestInput(t *testing.T) {
	var testData = []struct{ description, input, output string }{
		{
			description: "smoke test",
			input:       "",
			output:      "",
		},
	}

	for _, x := range testData {
		x := x
		t.Run(x.description, func(t *testing.T) {
			testMain(t, x.input, x.output)
		})
	}
}

func testMain(t *testing.T, input, expectedOutput string) {
	csvLog, err := ioutil.TempFile(os.TempDir(), "datadog-parser")
	assert.NoError(t, err)
	defer os.RemoveAll(csvLog.Name())

	_, err = csvLog.Write([]byte(input))
	assert.NoError(t, err)

	// prepare stdout capture
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"test", "--filepath=" + csvLog.Name()}
	done := make(chan struct{})
	go func() {
		<-done
		e := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		assert.NoError(t, e)
	}()
	finished := make(chan struct{})
	go func() {
		main()
		close(finished)
	}()

	// kill program after test is done
	defer func() {
		close(done)
		<-finished
	}()

	// awful hack to give program enough time to write output to stdout
	time.Sleep(time.Second)

	// restore stdout
	w.Close()
	os.Stdout = rescueStdout

	out, _ := ioutil.ReadAll(r)
	assert.Equal(t, expectedOutput, string(out))
}
