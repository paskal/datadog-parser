package record

import (
	"strconv"
	"strings"
	"time"
)

type record struct {
	remotehost string
	rfc931     string
	authuser   string
	date       time.Time
	request    string
	section    string
	status     int
	bytes      int
}

// parseRecord from slice of strings, return false in terms of errors
func parseRecord(raw []string) *record {
	if len(raw) != 7 {
		return nil
	}
	var err error
	r := record{
		remotehost: raw[0],
		rfc931:     raw[1],
		authuser:   raw[2],
		request:    raw[4],
	}
	if r.bytes, err = strconv.Atoi(raw[6]); err != nil {
		return nil
	}
	if r.status, err = strconv.Atoi(raw[5]); err != nil {
		return nil
	}
	var timestamp int64
	if timestamp, err = strconv.ParseInt(raw[3], 10, 64); err != nil {
		return nil
	}
	r.date = time.Unix(timestamp, 0)
	s := strings.Split(r.request, "/")
	if len(s) < 2 {
		return nil
	}
	r.section = s[1]
	return &r
}
