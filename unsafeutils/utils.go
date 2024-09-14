package unsafeutils

import (
	"reflect"
	"unsafe"
)

func Get[O any, T any](ptr *T, field string) (o O, ok bool) {
	mapping := FieldsOffset[T]()
	offset, ok := mapping[field]
	if !ok {
		return o, false
	}

	o = GetWithOffset[O](ptr, offset)
	return o, true
}

func Set[O any, T any](ptr *T, field string, i O) (ok bool) {
	mapping := FieldsOffset[T]()
	offset, ok := mapping[field]
	if !ok {
		return false
	}

	SetWithOffset[O](ptr, offset, i)
	return true
}

var cacheFields = map[reflect.Type][]string{}

func Fields[T any]() []string {
	t := reflect.TypeFor[T]()
	out, ok := cacheFields[t]
	if ok {
		return out
	}
	n := t.NumField()
	out = make([]string, 0, n)
	for i := 0; i < n; i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		out = append(out, field.Name)
	}
	cacheFields[t] = out
	return out
}

var cacheFieldsOffset = map[reflect.Type]map[string]uintptr{}

func FieldsOffset[T any]() map[string]uintptr {
	t := reflect.TypeFor[T]()
	out, ok := cacheFieldsOffset[t]
	if ok {
		return out
	}
	n := t.NumField()
	out = map[string]uintptr{}
	for i := 0; i < n; i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		out[field.Name] = field.Offset
	}
	cacheFieldsOffset[t] = out
	return out
}

func fieldWithOffset[O any, T any](ptr *T, offset uintptr) *O {
	return (*O)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + offset))
}

func GetWithOffset[O any, T any](ptr *T, offset uintptr) O {
	return *fieldWithOffset[O](ptr, offset)
}

func SetWithOffset[O any, T any](ptr *T, offset uintptr, i O) {
	*fieldWithOffset[O](ptr, offset) = i
}
