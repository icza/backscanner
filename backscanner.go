/*
Package backscanner provides a scanner similar to bufio.Scanner that reads
and returns lines in reverse order, starting at a given position (which may be
the end of the input).
*/
package backscanner

import (
	"bytes"
	"io"
)

// Scanner is the back-scanner implementation.
type Scanner struct {
	r   io.ReaderAt
	pos int
	err error
	buf []byte
}

// New returns a new Scanner.
func New(r io.ReaderAt, pos int) *Scanner {
	return &Scanner{r: r, pos: pos}
}

// readMore reads more data from the input.
func (s *Scanner) readMore() {
	if s.pos == 0 {
		s.err = io.EOF
		return
	}
	size := 1024
	if size > s.pos {
		size = s.pos
	}
	s.pos -= size
	buf2 := make([]byte, size, size+len(s.buf))

	// ReadAt attempts to read full buff!
	_, s.err = s.r.ReadAt(buf2, int64(s.pos))
	if s.err == nil {
		s.buf = append(buf2, s.buf...)
	}
}

// Line returns the next line from the input.
// Also returns the absolute start byte-position of the line in the input.
func (s *Scanner) Line() (line string, start int, err error) {
	if s.err != nil {
		return "", 0, s.err
	}
	for {
		lineStart := bytes.LastIndexByte(s.buf, '\n')
		if lineStart >= 0 {
			// We have a complete line:
			var line string
			line, s.buf = string(dropCR(s.buf[lineStart+1:])), s.buf[:lineStart]
			return line, s.pos + lineStart, nil
		}
		// Need more data:
		s.readMore()
		if s.err != nil {
			if s.err == io.EOF {
				if len(s.buf) > 0 {
					return string(dropCR(s.buf)), 0, nil
				}
			}
			return "", 0, s.err
		}
	}
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
