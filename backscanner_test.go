package backscanner

import (
	"io"
	"strings"
	"testing"

	"github.com/icza/mighty"
)

func TestDefaults(t *testing.T) {
	eq := mighty.Eq(t)

	scanner := New(nil, 0)
	eq(DefaultChunkSize, scanner.o.ChunkSize)
	eq(DefaultMaxBufferSize, scanner.o.MaxBufferSize)

	scanner = NewOptions(nil, 0, &Options{
		ChunkSize:     -1,
		MaxBufferSize: -1,
	})
	eq(DefaultChunkSize, scanner.o.ChunkSize)
	eq(DefaultMaxBufferSize, scanner.o.MaxBufferSize)
}

func TestScanner(t *testing.T) {
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
		// Leading empty lines must be returned too.
		{input: "\n", exps: []result{{"", 1, nil}, {"", 0, nil}, {"", 0, io.EOF}}},
		{input: "\nline", exps: []result{{"line", 1, nil}, {"", 0, nil}, {"", 0, io.EOF}}},
		{input: "\n\n", exps: []result{{"", 2, nil}, {"", 1, nil}, {"", 0, nil}, {"", 0, io.EOF}}},
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
		// Test with different chunk sizes:
		for _, chunkSize := range []int{-1, 0, 1, 2, 10, 100} {
			scanner := NewOptions(strings.NewReader(c.input), len(c.input), &Options{ChunkSize: chunkSize})
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
}

func TestLongLine(t *testing.T) {
	eq := mighty.Eq(t)

	scanner := NewOptions(strings.NewReader("123456789"), 9, &Options{
		MaxBufferSize: 5,
	})

	_, _, err := scanner.Line()
	eq(ErrLongLine, err)
}

type fullBufferAndEOFReaderAt struct {
	content string
}

func (r fullBufferAndEOFReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if len(p) == len(r.content) && off == 0 {
		copy(p, r.content)
		return len(p), io.EOF
	}
	return strings.NewReader(r.content).ReadAt(p, off)
}

func TestFullBufferAndEOF(t *testing.T) {
	eq := mighty.Eq(t)

	in := "1234567890"
	scanner := NewOptions(fullBufferAndEOFReaderAt{in}, len(in), &Options{
		MaxBufferSize: len(in),
	})

	line, pos, err := scanner.Line()
	eq(nil, err)
	eq(in, line)
	eq(0, pos)
}

func TestSmallBufferReadsShortLines(t *testing.T) {
	input := strings.Repeat("a\nb\n", 8) + "end"
	for _, chunkSize := range []int{0, 4, 5, 100} {
		scanner := NewOptions(strings.NewReader(input), len(input), &Options{
			ChunkSize: chunkSize, MaxBufferSize: 5,
		})
		parts := strings.Split(input, "\n")
		pos := len(input)
		for i := len(parts) - 1; i >= 0; i-- {
			pos -= len(parts[i])
			line, gotPos, err := scanner.Line()
			if err != nil || line != parts[i] || gotPos != pos {
				t.Fatalf("chunk=%d line=%d: got (%q, %d, %v), want (%q, %d, nil)", chunkSize, i, line, gotPos, err, parts[i], pos)
			}
			pos--
		}
		if _, _, err := scanner.Line(); err != io.EOF {
			t.Fatalf("chunk=%d: want EOF, got %v", chunkSize, err)
		}
	}
}
