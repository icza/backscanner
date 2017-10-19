# backscanner

[![Build Status](https://travis-ci.org/icza/backscanner.svg?branch=master)](https://travis-ci.org/icza/backscanner)
[![GoDoc](https://godoc.org/github.com/icza/backscanner?status.svg)](https://godoc.org/github.com/icza/backscanner)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/backscanner)](https://goreportcard.com/report/github.com/icza/backscanner)
[![codecov](https://codecov.io/gh/icza/backscanner/branch/master/graph/badge.svg)](https://codecov.io/gh/icza/backscanner)

Ever needed or wondered how to efficiently search for something in a log file,
but starting at the end and going backward? Here's your solution.

Package `backscanner` provides a scanner similar to `bufio.Scanner`, but it reads
and returns lines in reverse order, starting at a given position (which may be
the end of the input) and going backward.

Advancing and accessing lines of the input is done by calling the `Scanner.Line()`
method, which returns the next line (previous in the source) as a `string`.

For maximum efficiency, there is a `Scanner.LineBytes()` method. This method
returns the next line as a byte slice, which shares its backing array with the
internal buffer of the `Scanner`. This is for efficiency so no copy is made from
the line data, but this also means you can only inspect or search in the slice
before calling `Line()` or `LineBytes()` again, as the content of the internal
buffer–and thus slices returned by `LineBytes()`–may be overwritten. If you need
to retain the line data, make a copy of it or use the `Line()` method.


Example using it:

	input := "Line1\nLine2\nLine3"
	scanner := backscanner.New(strings.NewReader(input), len(input))
	for {
		line, pos, err := scanner.Line()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		fmt.Printf("Line position: %2d, line: %q\n", pos, line)
	}

Output:

	Line position: 12, line: "Line3"
	Line position:  6, line: "Line2"
	Line position:  0, line: "Line1"
	Error: EOF

Using it to scan a file backward, starting from its end:

	f, err := os.Open("mylog.txt")
	if err != nil {
		panic(err)
	}
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := backscanner.New(f, int(fi.Size()))
	// Now use scanner like in the previous example
