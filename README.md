# backscanner

[![Build Status](https://travis-ci.org/icza/backscanner.svg?branch=master)](https://travis-ci.org/icza/backscanner)
[![GoDoc](https://godoc.org/github.com/icza/backscanner?status.svg)](https://godoc.org/github.com/icza/backscanner)
[![Go Report Card](https://goreportcard.com/badge/github.com/icza/backscanner)](https://goreportcard.com/report/github.com/icza/backscanner)

Package `backscanner` provides a scanner similar to `bufio.Scanner` that reads
and returns lines in reverse order, starting at a given position (which may be
the end of the input).
