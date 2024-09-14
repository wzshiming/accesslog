package accesslog

import (
	"fmt"
	"strings"

	"github.com/wzshiming/accesslog/unsafeutils"
)

type Entry[T any] struct {
	entry T
}

func ParseEntry[T any](raw []byte) (entry Entry[T], err error) {
	data := raw
	fields := entry.Fields()
	for _, field := range fields {
		cur, next, ok := parseField(data)
		if !ok {
			return entry, fmt.Errorf("failed to parse field %q: %q", field, raw)
		}
		if cur == nil {
			break
		}
		ok = unsafeutils.Set(&entry.entry, field, string(cur))
		if !ok {
			return entry, fmt.Errorf("failed to set field %q: %q", field, raw)
		}
		data = next
		if len(data) == 0 {
			break
		}
	}
	if len(data) != 0 {
		return entry, fmt.Errorf("unexpected data remaining: %q", data)
	}
	return entry, nil
}

func (e *Entry[T]) Entry() *T {
	return &e.entry
}

func (*Entry[T]) Fields() []string {
	return unsafeutils.Fields[T]()
}

func (e *Entry[T]) Values(fields []string) []string {
	accessLogEntryFieldsIndexMapping := unsafeutils.FieldsOffset[T]()
	out := make([]string, len(fields))
	for i, f := range fields {
		offset, ok := accessLogEntryFieldsIndexMapping[f]
		if !ok {
			continue
		}
		out[i] = unsafeutils.GetWithOffset[string](&e.entry, offset)
	}
	return out
}

func (e *Entry[T]) String() string {
	return strings.Join(e.Values(e.Fields()), " ")
}
