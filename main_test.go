package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	r := strings.NewReader(rnString(`BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:aa
DESCRIPTION:aacc
END:VEVENT
BEGIN:VEVENT
SUMMARY:aaabbb
DESCRIPTION:aacc
END:VEVENT
BEGIN:VEVENT
SUMMARY:bb
DESCRIPTION:ddd
END:VEVENT
BEGIN:VEVENT
SUMMARY:111
DESCRIPTION:aaaccc
END:VEVENT
END:VCALENDAR
`))
	var b bytes.Buffer
	if err := _main(r, &b, []string{`SUMMARY:aaa.*`, `SUMMARY:bb`}, []string{`^DESCRIPTION:.*ccc$`}); err != nil {
		t.Fatal(err)
	}
	if e, g := rnString(`BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
SUMMARY:aa
DESCRIPTION:aacc
END:VEVENT
BEGIN:VEVENT
SUMMARY:111
END:VEVENT
END:VCALENDAR
`), b.String(); e != g {
		t.Errorf("expected %v got %v", e, g)
	}
}

func rnString(s string) string {
	return strings.Replace(s, "\n", "\r\n", -1)
}
