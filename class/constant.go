package class

import (
	"fmt"
)

const (
	_class              = 7
	_fieldRef           = 9
	_methodRef          = 10
	_interfaceMethodRef = 11
	_string             = 8
	_integer            = 3
	_float              = 4
	_long               = 5
	_double             = 6
	_nameAndType        = 12
	_utf8               = 1
	_methodHandle       = 15
	_methodType         = 16
	_invokeDynamic      = 18
)

var (
	_tm = map[uint8]string{
		_class:              "Class",
		_fieldRef:           "Fieldref",
		_methodRef:          "Methodref",
		_interfaceMethodRef: "Methodref",
		_string:             "String",
		_integer:            "Integer",
		_float:              "Float",
		_long:               "Long",
		_double:             "Double",
		_nameAndType:        "NameAndType",
		_utf8:               "utf8",
		_methodHandle:       "MethodHandle",
		_methodType:         "MethodType",
		_invokeDynamic:      "InvokeDynamic",
	}
)

func NewConstantInfo(b []byte, s int) (res ConstantInfo, next int) {
	tag, next := u8(b, s)
	switch tag {
	case _class:
		res = new(ClassInfo)
	case _fieldRef:
		res = new(FieldRefInfo)
	case _methodRef:
		res = new(MethodRefInfo)
	case _interfaceMethodRef:
		res = new(InterfaceMethodRefInfo)
	case _string:
		res = new(StringInfo)
	case _integer:
		res = new(IntegerInfo)
	case _float:
		res = new(FloatInfo)
	case _long:
		res = new(LongInfo)
	case _double:
		res = new(DoubleInfo)
	case _nameAndType:
		res = new(NameAndType)
	case _utf8:
		res = new(Utf8Info)
	case _methodHandle:
		res = new(MethodHandle)
	case _methodType:
		res = new(MethodTypeInfo)
	case _invokeDynamic:
		res = new(InvokeDynamicInfo)
	default:
		panic(fmt.Errorf("unsupported tag %d", tag))
	}
	res.SetT(tag)
	return res, res.Read(b, next)
}

// ConstantInfo constant info, each constant info holds tag.
type ConstantInfo interface {
	T() uint8
	TN() string
	SetT(tag uint8)
	Read(b []byte, s int) (next int)
}

type Tag struct {
	Tag uint8
}

func (m *Tag) T() uint8 {
	return m.Tag
}

func (m *Tag) TN() string {
	return _tm[m.Tag]
}

func (m *Tag) SetT(tag uint8) {
	m.Tag = tag
}

// ClassInfo CONSTANT_Class_info.
type ClassInfo struct {
	Tag
	NameIndex uint16
}

func (m *ClassInfo) Read(b []byte, s int) (next int) {
	m.NameIndex, next = u16(b, s)
	return
}

func (m *ClassInfo) ParseNameFromPool(cp []ConstantInfo) (name string, err error) {
	u2, ok := cp[m.NameIndex].(*Utf8Info)
	if !ok {
		err = fmt.Errorf("NameIndex of class info point to a non utf8 info")
		return
	}
	name, err = u2string(u2.Bytes)
	return
}

// FiledRefInfo CONSTANT_Fieldref_info.
type FieldRefInfo struct {
	Tag
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

func (m *FieldRefInfo) Read(b []byte, s int) (next int) {
	m.ClassIndex, next = u16(b, s)
	m.NameAndTypeIndex, next = u16(b, next)
	return
}

func (m *FieldRefInfo) ParseClassFromPool(cp []ConstantInfo) (class string, err error) {
	n, ok := cp[m.ClassIndex].(*ClassInfo)
	if !ok {
		err = fmt.Errorf("class index must pointer to a ClassInfo")
		return
	}
	return n.ParseNameFromPool(cp)
}

func (m *FieldRefInfo) ParseNameAndTypeFromPool(cp []ConstantInfo) (name, desc string, err error) {
	n, ok := cp[m.NameAndTypeIndex].(*NameAndType)
	if !ok {
		err = fmt.Errorf("name and type index must pointer to a NameAndType")
		return
	}
	name, desc, err = n.ParseFromPool(cp)
	return
}

// MethodRefInfo CONSTANT_Methodref_info.
type MethodRefInfo struct {
	FieldRefInfo
}

// InterfaceMethodRefInfo CONSTANT_InterfaceMethodref_info.
type InterfaceMethodRefInfo struct {
	FieldRefInfo
}

// StringInfo CONSTANT_String_info.
type StringInfo struct {
	Tag
	StringIndex uint16
}

func (m *StringInfo) ParseStringFromPool(cp []ConstantInfo) (s string, err error) {
	return ui2string(cp, m.StringIndex)
}

func (m *StringInfo) Read(b []byte, s int) (next int) {
	m.StringIndex, next = u16(b, s)
	return
}

// IntegerInfo CONSTANT_Integer_info.
type IntegerInfo struct {
	Tag
	Bytes uint32
}

func (m *IntegerInfo) Read(b []byte, s int) (next int) {
	m.Bytes, next = u32(b, s)
	return
}

// FloatInfo CONSTANT_Float_info.
type FloatInfo struct {
	IntegerInfo
}

// LongInfo CONSTANT_Long_info.
type LongInfo struct {
	Tag
	HighBytes uint32
	LowBytes  uint32
}

func (m *LongInfo) Read(b []byte, s int) (next int) {
	m.HighBytes, next = u32(b, s)
	m.LowBytes, next = u32(b, next)
	return
}

// DoubleInfo CONSTANT_Double_info.
type DoubleInfo struct {
	LongInfo
}

// NameAndTypeInfo CONSTANT_NameAndType_info.
type NameAndType struct {
	Tag
	NameIndex       uint16
	DescriptorIndex uint16
}

func (m *NameAndType) Read(b []byte, s int) (next int) {
	m.NameIndex, next = u16(b, s)
	m.DescriptorIndex, next = u16(b, next)
	return
}

func (m *NameAndType) ParseFromPool(cp []ConstantInfo) (name, desc string, err error) {
	if name, err = m.ParseNameFromPool(cp); err != nil {
		return
	}
	desc, err = m.ParseDescFromPool(cp)
	return
}

func (m *NameAndType) ParseNameFromPool(cp []ConstantInfo) (res string, err error) {
	return ui2string(cp, m.NameIndex)
}

func (m *NameAndType) ParseDescFromPool(cp []ConstantInfo) (res string, err error) {
	return ui2string(cp, m.DescriptorIndex)
}

func ui2string(cp []ConstantInfo, i uint16) (res string, err error) {
	u2, ok := cp[i].(*Utf8Info)
	if !ok {
		err = fmt.Errorf("index %d points to a non utf8 info", i)
		return
	}
	res, err = u2string(u2.Bytes)
	return
}

// Utf8Info CONSTANT_Utf8_info.
type Utf8Info struct {
	Tag
	Length uint16
	Bytes  []byte
}

func (m *Utf8Info) Read(b []byte, s int) (next int) {
	m.Length, next = u16(b, s)
	m.Bytes, next = bs(b, next, int(m.Length))
	return
}

// MethodHandle CONSTANT_MethodHandle_info.
type MethodHandle struct {
	Tag
	ReferenceKind  uint8
	ReferenceIndex uint16
}

func (m *MethodHandle) Read(b []byte, s int) (next int) {
	m.ReferenceKind, next = u8(b, s)
	m.ReferenceIndex, next = u16(b, next)
	return
}

// MethodTypeInfo CONSTANT_MethodType_info.
type MethodTypeInfo struct {
	Tag
	DescriptorIndex uint16
}

func (m *MethodTypeInfo) Read(b []byte, s int) (next int) {
	m.DescriptorIndex, next = u16(b, s)
	return
}

// InvokeDynamicInfo CONSTANT_InvokeDynamic_info.
type InvokeDynamicInfo struct {
	Tag
	BootstrapMethodAttrIndex uint16
	NameAndTypeIndex         uint16
}

func (m *InvokeDynamicInfo) Read(b []byte, s int) (next int) {
	m.BootstrapMethodAttrIndex, next = u16(b, s)
	m.NameAndTypeIndex, next = u16(b, next)
	return
}
