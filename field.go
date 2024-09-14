package accesslog

import (
	"bytes"
)

func parseField(data []byte) (cur []byte, next []byte, ok bool) {
	data = bytes.TrimLeft(data, " ")
	if len(data) == 0 {
		return nil, nil, true
	}

	switch data[0] {
	case '"':
		index := bytes.IndexByte(data[1:], '"')
		if index == -1 {
			return nil, nil, false
		}
		return data[:index+2], data[index+2:], true
	case '[':
		index := bytes.IndexByte(data[1:], ']')
		if index == -1 {
			return nil, nil, false
		}
		return data[:index+2], data[index+2:], true
	}

	index := bytes.IndexByte(data, ' ')
	if index == -1 {
		return data, nil, true
	}
	return data[:index], data[index+1:], true
}
