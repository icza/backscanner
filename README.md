# backscanner

[![Build Status](https://travis-ci.org/icza/backscanner.svg?branch=master)](https://travis-ci.org/icza/backscanner)
[![GoDoc](https://godoc.org/github.com/icza/backscanner?status.svg)](https://godoc.org/github.com/icza/backscanner)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/backscanner)](https://goreportcard.com/report/github.com/icza/backscanner)

Ever needed or wondered how to efficiently search for something in a log file,
but starting at the end and going backward? Here's a solution now.

Package `backscanner` provides a scanner similar to `bufio.Scanner` that reads
and returns lines in reverse order, starting at a given position (which may be
the end of the input).

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
