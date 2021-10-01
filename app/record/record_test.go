package record

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseRecord(t *testing.T) {
	var testData = []struct {
		input  []string
		output *record
	}{
		{input: []string{}, output: nil},
		{input: []string{"remotehost", "rfc931", "authuser", "date", "request", "status", "bytes"}, output: nil},
		{input: []string{"10.0.0.2", "-", "apache", "not_a_number", "GET /api/user HTTP/1.0", "200", "1234"}, output: nil},
		{input: []string{"10.0.0.2", "-", "apache", "1549573860", "GET /api/user HTTP/1.0", "not_a_number", "1234"}, output: nil},
		{input: []string{"10.0.0.2", "-", "apache", "1549573860", "GET /api/user HTTP/1.0", "200", "not_a_number"}, output: nil},
		{input: []string{"10.0.0.2", "-", "apache", "1549573860", "wrong_line", "200", "1234"}, output: nil},
		{input: []string{"10.0.0.2", "-", "apache", "1549573860", "GET /api/user HTTP/1.0", "200", "1234", "extra_field"}, output: nil},
		{
			input: []string{"10.0.0.2", "-", "apache", "1549573860", "GET /api/user HTTP/1.0", "200", "1234"},
			output: &record{
				remotehost: "10.0.0.2",
				rfc931:     "-",
				authuser:   "apache",
				date:       time.Date(2019, 02, 07, 21, 11, 0, 0, time.UTC).In(time.Local),
				request:    "GET /api/user HTTP/1.0",
				section:    "api",
				status:     200,
				bytes:      1234,
			},
		},
	}

	for i, x := range testData {
		x := x
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, x.output, parseRecord(x.input))
		})
	}
}
