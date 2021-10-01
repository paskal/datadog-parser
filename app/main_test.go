package main

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const sampleCsvOutput = `2019-02-07 21:11:09 +0000 UTC: 81 hits from 5 users with 99752 bytes transferred, top /api with 11 hits
2019-02-07 21:11:19 +0000 UTC: 94 hits from 5 users with 116359 bytes transferred, top /api with 11 hits
2019-02-07 21:11:29 +0000 UTC: 99 hits from 5 users with 121644 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:11:39 +0000 UTC: 100 hits from 5 users with 122898 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:11:49 +0000 UTC: 93 hits from 5 users with 113623 bytes transferred, top /api with 11 hits
2019-02-07 21:11:59 +0000 UTC: 92 hits from 5 users with 112364 bytes transferred, top /api with 11 hits
2019-02-07 21:12:09 +0000 UTC: 171 hits from 5 users with 210066 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:12:19 +0000 UTC: 181 hits from 5 users with 222160 bytes transferred, top /api with 11 hits
2019-02-07 21:12:29 +0000 UTC: 182 hits from 5 users with 223896 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:12:36 +0000 UTC: Alert RED, ~10.01 hits per second which is higher than 10 (1201 total) in the last 2m0s
2019-02-07 21:12:39 +0000 UTC: 191 hits from 5 users with 235645 bytes transferred, top /api with 11 hits
2019-02-07 21:12:49 +0000 UTC: 178 hits from 5 users with 219456 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:12:59 +0000 UTC: 190 hits from 5 users with 233590 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:13:09 +0000 UTC: 49 hits from 5 users with 60201 bytes transferred, top /api with 10 hits
2019-02-07 21:13:19 +0000 UTC: 33 hits from 5 users with 41284 bytes transferred, top /api and /report with 10 hits
2019-02-07 21:13:29 +0000 UTC: 33 hits from 5 users with 41217 bytes transferred, top /api with 10 hits
2019-02-07 21:13:39 +0000 UTC: 32 hits from 5 users with 39241 bytes transferred, top /api with 11 hits
2019-02-07 21:13:49 +0000 UTC: 32 hits from 5 users with 39501 bytes transferred, top /api with 11 hits
2019-02-07 21:13:59 +0000 UTC: 33 hits from 5 users with 40646 bytes transferred, top /api with 10 hits
2019-02-07 21:14:04 +0000 UTC: Alert GREEN, ~9.97 hits per second which is lower than 10 (1197 total) in the last 2m0s
2019-02-07 21:14:04 +0000 UTC: Alert RED, ~10.01 hits per second which is higher than 10 (1201 total) in the last 2m0s
2019-02-07 21:14:05 +0000 UTC: Alert GREEN, ~9.89 hits per second which is lower than 10 (1187 total) in the last 2m0s
2019-02-07 21:14:09 +0000 UTC: 30 hits from 5 users with 37185 bytes transferred, top /api with 10 hits
2019-02-07 21:14:18 +0000 UTC: 32 hits from 5 users with 38818 bytes transferred, top /api and /report with 10 hits
2019-02-07 21:14:29 +0000 UTC: 32 hits from 5 users with 39059 bytes transferred, top /api with 11 hits
2019-02-07 21:14:39 +0000 UTC: 34 hits from 5 users with 41862 bytes transferred, top /api with 11 hits
2019-02-07 21:14:49 +0000 UTC: 34 hits from 5 users with 41614 bytes transferred, top /api with 11 hits
2019-02-07 21:14:59 +0000 UTC: 35 hits from 5 users with 42452 bytes transferred, top /api with 10 hits
2019-02-07 21:15:09 +0000 UTC: 34 hits from 5 users with 41539 bytes transferred, top /api with 11 hits
2019-02-07 21:15:19 +0000 UTC: 33 hits from 5 users with 40386 bytes transferred, top /api with 10 hits
2019-02-07 21:15:29 +0000 UTC: 33 hits from 5 users with 40601 bytes transferred, top /report with 11 hits
2019-02-07 21:15:39 +0000 UTC: 257 hits from 5 users with 316992 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:15:49 +0000 UTC: 287 hits from 5 users with 353279 bytes transferred, top /api with 11 hits
2019-02-07 21:15:59 +0000 UTC: 283 hits from 5 users with 348392 bytes transferred, top /api with 11 hits
2019-02-07 21:16:03 +0000 UTC: Alert RED, ~10.01 hits per second which is higher than 10 (1201 total) in the last 2m0s
2019-02-07 21:16:09 +0000 UTC: 283 hits from 5 users with 347894 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:16:19 +0000 UTC: 284 hits from 5 users with 348956 bytes transferred, top /api with 11 hits
2019-02-07 21:16:29 +0000 UTC: 280 hits from 5 users with 344621 bytes transferred, top /api with 11 hits
2019-02-07 21:16:39 +0000 UTC: 282 hits from 5 users with 346313 bytes transferred, top /api and /report with 11 hits
2019-02-07 21:16:49 +0000 UTC: 303 hits from 5 users with 372471 bytes transferred, top /api with 11 hits
2019-02-07 21:16:59 +0000 UTC: 279 hits from 5 users with 342950 bytes transferred, top /api with 11 hits
2019-02-07 21:17:09 +0000 UTC: 46 hits from 5 users with 56366 bytes transferred, top /api and /report with 9 hits
2019-02-07 21:17:20 +0000 UTC: 22 hits from 5 users with 26888 bytes transferred, top /api and /report with 8 hits
2019-02-07 21:17:30 +0000 UTC: 21 hits from 5 users with 25581 bytes transferred, top /api with 9 hits
2019-02-07 21:17:40 +0000 UTC: 23 hits from 5 users with 28189 bytes transferred, top /report with 9 hits
2019-02-07 21:17:50 +0000 UTC: 23 hits from 5 users with 28546 bytes transferred, top /report with 8 hits
2019-02-07 21:18:00 +0000 UTC: 21 hits from 5 users with 25952 bytes transferred, top /report with 8 hits
2019-02-07 21:18:10 +0000 UTC: 21 hits from 4 users with 25905 bytes transferred, top /api with 9 hits
2019-02-07 21:18:20 +0000 UTC: 21 hits from 5 users with 25635 bytes transferred, top /api with 9 hits
2019-02-07 21:18:23 +0000 UTC: Alert GREEN, ~9.97 hits per second which is lower than 10 (1197 total) in the last 2m0s
2019-02-07 21:18:30 +0000 UTC: 20 hits from 5 users with 24603 bytes transferred, top /api with 8 hits
2019-02-07 21:18:40 +0000 UTC: 23 hits from 5 users with 28184 bytes transferred, top /api with 10 hits
2019-02-07 21:18:50 +0000 UTC: 22 hits from 5 users with 26888 bytes transferred, top /api with 11 hits
`

func TestSampleCSV(t *testing.T) {
	testMain(t, "../sample.csv", sampleCsvOutput)
}

func TestInput(t *testing.T) {
	var testData = []struct{ description, input, output string }{
		{
			description: "smoke test",
			input:       "",
			output:      "",
		},
		{
			description: "single line, no output expected",
			input: `"remotehost","rfc931","authuser","date","request","status","bytes"
"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234
`,
			output: "",
		},
		{
			description: "two lines, previous stats should be reported",
			input: `"remotehost","rfc931","authuser","date","request","status","bytes"
"10.0.0.2","-","apache",1549573860,"GET /api/user HTTP/1.0",200,1234
"10.0.0.1","-","apache",1549573891,"POST /report HTTP/1.0",500,1234
`,
			output: `2019-02-07 21:11:00 +0000 UTC: 1 hits from 1 users with 1234 bytes transferred, top /api with 1 hits
`,
		},
	}

	for _, x := range testData {
		x := x
		t.Run(x.description, func(t *testing.T) {
			csvLog, err := ioutil.TempFile(os.TempDir(), "datadog-parser")
			assert.NoError(t, err)
			defer os.RemoveAll(csvLog.Name())

			_, err = csvLog.Write([]byte(x.input))
			assert.NoError(t, err)

			testMain(t, csvLog.Name(), x.output)
		})
	}
}

func testMain(t *testing.T, inputFile, expectedOutput string) {
	// prepare stdout capture
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"test", "--filepath=" + inputFile}
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
