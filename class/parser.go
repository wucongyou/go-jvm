package class

import (
	"encoding/binary"
	"io/ioutil"
)

func ParseFile(filename string) (res *ClassFile, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return ParseBytes(b)
}

func ParseBytes(b []byte) (res *ClassFile, err error) {
	res = new(ClassFile)
	var next int
	res.Magic, next = u32(b, 0)
	res.MinorVersion, next = u16(b, next)
	res.MajorVersion, next = u16(b, next)

	res.ConstantPoolCount, next = u16(b, next)
	// parse constant pool: #1 - #constant_pool_count-1
	res.CpInfo = make([]ConstantInfo, res.ConstantPoolCount)
	var c ConstantInfo
	for i := 1; i < int(res.ConstantPoolCount); i++ {
		c, next = NewConstantInfo(b, next)
		res.CpInfo[i] = c
		if c.T() == _double || c.T() == _long {
			i++
		}
	}

	res.AccessFlags, next = u16(b, next)
	res.ThisClass, next = u16(b, next)
	res.SuperClass, next = u16(b, next)

	res.InterfacesCount, next = u16(b, next)
	res.Interfaces = make([]*ClassInfo, res.InterfacesCount)
	for i := 0; i < int(res.InterfacesCount); i++ {
		res.Interfaces[i] = new(ClassInfo)
		next = res.Interfaces[i].Read(b, next)
	}

	res.FieldsCount, next = u16(b, next)
	res.Fields = make([]*FieldInfo, res.FieldsCount)
	for i := 0; i < int(res.FieldsCount); i++ {
		res.Fields[i] = new(FieldInfo)
		next = res.Fields[i].Read(b, next)
	}

	res.MethodsCount, next = u16(b, next)
	res.Methods = make([]*MethodInfo, res.MethodsCount)
	for i := 0; i < int(res.MethodsCount); i++ {
		res.Methods[i] = new(MethodInfo)
		next = res.Methods[i].Read(b, next)
	}

	res.AttributesCount, next = u16(b, next)
	res.Attributes = make([]*AttributeInfo, res.AttributesCount)
	for i := 0; i < int(res.AttributesCount); i++ {
		res.Attributes[i] = new(AttributeInfo)
		next = res.Attributes[i].Read(b, next)
	}
	return
}

func u8(b []byte, s int) (res uint8, next int) {
	next = s + 1
	res = uint8(b[s])
	return
}

func u16(b []byte, s int) (res uint16, next int) {
	next = s + 2
	res = binary.BigEndian.Uint16(b[s:next])
	return
}

func u32(b []byte, s int) (res uint32, next int) {
	next = s + 4
	res = binary.BigEndian.Uint32(b[s:next])
	return
}

func bs(b []byte, s int, len int) (res []byte, next int) {
	next = s + len
	res = b[s:next]
	return
}
