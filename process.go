package accesslog

import (
	"bufio"
	"io"
	"sync"
)

var bufReaderPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReader(nil)
	},
}

func ProcessEntries[T any](r io.Reader, callback func(entry Entry[T], err error) error) error {
	scanner := bufReaderPool.Get().(*bufio.Reader)
	scanner.Reset(r)
	defer func() {
		scanner.Reset(nil)
		bufReaderPool.Put(scanner)
	}()

	for {
		line, _, err := scanner.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = callback(ParseEntry[T](line))
		if err != nil {
			return err
		}
	}
	return nil
}
