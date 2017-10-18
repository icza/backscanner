package backscanner

import (
	"io"
	"strings"
	"testing"

	"github.com/icza/mighty"
)

func Test1(t *testing.T) {
	eq := mighty.Eq(t)

	type result struct {
		line string
		pos  int
		err  error
	}

	cases := []struct {
		input string
		exps  []result
	}{
		// Empty input
		{input: "", exps: []result{{"", 0, io.EOF}}},
		// Normal input with \n line endings
		{
			input: "Start\nLine1\nLine2\nLine3\nEnd",
			exps: []result{
				{"End", 24, nil},
				{"Line3", 18, nil},
				{"Line2", 12, nil},
				{"Line1", 6, nil},
				{"Start", 0, nil},
				{"", 0, io.EOF},
			},
		},
		// Normal input with \r\n line endings
		{
			input: "Line1\r\nLine2\r\n",
			exps: []result{
				{"", 14, nil},
				{"Line2", 7, nil},
				{"Line1", 0, nil},
				{"", 0, io.EOF},
			},
		},
	}

	for _, c := range cases {
		scanner := New(strings.NewReader(c.input), len(c.input))
		i := 0
		for {
			line, pos, err := scanner.Line()
			exp := c.exps[i]
			eq(exp.line, line)
			eq(exp.pos, pos)
			eq(exp.err, err)
			if err == io.EOF {
				eq(len(c.exps)-1, i)
				break
			}
			i++
		}
	}
}
