package calendar

import (
	"bufio"
	"io"
	"strings"
)

// Scanner reads a io.Reader line by line and tokenise it by key/value.
type Scanner struct {
	scanner *bufio.Scanner
}

// Scan parses the next line from the io.Reader.
func (r *Scanner) Scan() bool {
	for r.scanner.Scan() {
		if strings.Index(r.scanner.Text(), ":") != -1 {
			return true
		}
	}

	return false
}

// KeyValue returns the key/value present on the current line.
func (r *Scanner) KeyValue() (key string, value string) {
	line := r.scanner.Text()

	splits := strings.SplitN(line, ":", 2)
	switch len(splits) {
	case 1:
		key = splits[0]
		value = ""
		return
	case 2:
		key = splits[0]
		value = splits[1]
		return
	case 0:
		fallthrough
	default:
		panic("unreachable")
	}

	return
}

func (r *Scanner) Err() error {
	return r.scanner.Err()
}

func NewScanner(reader io.Reader) (scanner Scanner) {
	scanner.scanner = bufio.NewScanner(reader)
	return
}
