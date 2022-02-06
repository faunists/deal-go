package mocks

import "google.golang.org/protobuf/reflect/protoreflect"

var _ protoreflect.List = &ProtoList{}

type ProtoList struct {
	Values []protoreflect.Value
}

func (l *ProtoList) Append(v protoreflect.Value) {
	panic("unimplemented")
}

func (l *ProtoList) AppendMutable() protoreflect.Value {
	panic("unimplemented")
}

func (l *ProtoList) Get(n int) protoreflect.Value {
	return l.Values[n]
}

func (l *ProtoList) Len() int {
	return len(l.Values)
}

func (l *ProtoList) Set(n int, v protoreflect.Value) {
	panic("unimplemented")
}

func (l *ProtoList) Truncate(n int) {
	panic("unimplemented")
}

func (l *ProtoList) NewElement() protoreflect.Value {
	panic("unimplemented")
}

func (l *ProtoList) IsValid() bool {
	return true
}
