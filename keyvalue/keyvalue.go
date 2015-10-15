package keyvalue

import (
	"bufio"
	"os"
	"strings"
)

type Reader struct {
	scanner *bufio.Scanner
}

func (r *Reader) Next() (key string, value string, err error) {
	for r.scanner.Scan() {
		line := r.scanner.Text()
		err = r.scanner.Err()

		splits := strings.SplitN(line, ":", 2)
		switch len(splits) {
		case 0:
			continue
		case 1:
			key = splits[0]
			value = ""
			return
		case 2:
			key = splits[0]
			value = splits[1]
			return
		default:
			panic("unreachable")
		}
	}

	return
}

func NewReader(file *os.File) (reader Reader) {
	reader.scanner = bufio.NewScanner(file)
	return
}
