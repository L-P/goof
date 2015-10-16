package calendar

import (
	"bufio"
	"os"
	"strings"
)

type Scanner struct {
	scanner *bufio.Scanner
}

func (r *Scanner) Scan() bool {
	for r.scanner.Scan() {
		if strings.Index(r.scanner.Text(), ":") != -1 {
			return true
		}
	}

	return false
}

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

func NewScanner(file *os.File) (reader Scanner) {
	reader.scanner = bufio.NewScanner(file)
	return
}
