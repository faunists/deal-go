package mocks

import "google.golang.org/protobuf/reflect/protoreflect"

var _ protoreflect.Map = &ProtoMap{}

type ProtoMap struct {
	Map map[interface{}]protoreflect.Value
}

func (p *ProtoMap) Len() int {
	return len(p.Map)
}

func (p *ProtoMap) Range(f func(protoreflect.MapKey, protoreflect.Value) bool) {
	for key, value := range p.Map {
		if !f(protoreflect.ValueOf(key).MapKey(), value) {
			return
		}
	}
}

func (p *ProtoMap) Has(key protoreflect.MapKey) bool {
	panic("implement me")
}

func (p *ProtoMap) Clear(key protoreflect.MapKey) {
	panic("implement me")
}

func (p *ProtoMap) Get(key protoreflect.MapKey) protoreflect.Value {
	panic("implement me")
}

func (p *ProtoMap) Set(key protoreflect.MapKey, value protoreflect.Value) {
	panic("implement me")
}

func (p *ProtoMap) Mutable(key protoreflect.MapKey) protoreflect.Value {
	panic("implement me")
}

func (p *ProtoMap) NewValue() protoreflect.Value {
	panic("implement me")
}

func (p *ProtoMap) IsValid() bool {
	return true
}
