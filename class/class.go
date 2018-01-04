package class

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// ClassFile class file.
type ClassFile struct {
	Magic             uint32
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16
	CpInfo            []ConstantInfo
	AccessFlags       uint16
	ThisClass         uint16
	SuperClass        uint16
	InterfacesCount   uint16
	Interfaces        []*ClassInfo
	FieldsCount       uint16
	Fields            []*FieldInfo
	MethodsCount      uint16
	Methods           []*MethodInfo
	AttributesCount   uint16
	Attributes        []*AttributeInfo
}

func (m *ClassFile) Format() (res string, err error) {
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("magic: %x\n", m.Magic))
	b.WriteString(fmt.Sprintf("minor version: %d\n", m.MinorVersion))
	b.WriteString(fmt.Sprintf("major version: %d\n", m.MajorVersion))

	b.WriteString(fmt.Sprintf("constant pool count: %d\n", m.ConstantPoolCount))
	b.WriteString("constant pool:\n")
	l := 19
	c := len(strconv.FormatInt(int64(m.ConstantPoolCount), 10))
	for i := 1; i < int(m.ConstantPoolCount); i++ {
		if m.CpInfo[i] == nil {
			continue
		}
		c2 := len(strconv.FormatInt(int64(i), 10))
		b.WriteString(strings.Repeat(" ", c-c2))
		b.WriteString(fmt.Sprintf("#%d = ", i))
		c := m.CpInfo[i]
		n := c.TN()
		b.WriteString(n)
		b.WriteString(strings.Repeat(" ", l-len(n)))
		cm := ""
		cmShow := false
		switch u := m.CpInfo[i].(type) {
		case *ClassInfo:
			s := fmt.Sprintf("#%d", u.NameIndex)
			b.WriteString(s)
			b.WriteString(strings.Repeat(" ", l-len(s)))
			if cm, err = u.ParseNameFromPool(m.CpInfo); err != nil {
				return
			}

		case *FieldRefInfo:
			s := fmt.Sprintf("#%d:#%d", u.ClassIndex, u.NameAndTypeIndex)
			b.WriteString(s)
			b.WriteString(strings.Repeat(" ", l-len(s)))
			var class, name, desc string
			if class, err = u.ParseClassFromPool(m.CpInfo); err != nil {
				return
			}
			if name, desc, err = u.ParseNameAndTypeFromPool(m.CpInfo); err != nil {
				return
			}
			cm = fmt.Sprintf("%s.%s:%s", class, name, desc)
		case *MethodRefInfo:
			s := fmt.Sprintf("#%d:#%d", u.ClassIndex, u.NameAndTypeIndex)
			b.WriteString(s)
			b.WriteString(strings.Repeat(" ", l-len(s)))
			var class, name, desc string
			if class, err = u.ParseClassFromPool(m.CpInfo); err != nil {
				return
			}
			if name, desc, err = u.ParseNameAndTypeFromPool(m.CpInfo); err != nil {
				return
			}
			cm = fmt.Sprintf("%s.%s:%s", class, name, desc)
		case *InterfaceMethodRefInfo:
			s := fmt.Sprintf("#%d:#%d", u.ClassIndex, u.NameAndTypeIndex)
			b.WriteString(s)
			b.WriteString(strings.Repeat(" ", l-len(s)))
			var class, name, desc string
			if class, err = u.ParseClassFromPool(m.CpInfo); err != nil {
				return
			}
			if name, desc, err = u.ParseNameAndTypeFromPool(m.CpInfo); err != nil {
				return
			}
			cm = fmt.Sprintf("%s.%s:%s", class, name, desc)

		case *StringInfo:
			s := fmt.Sprintf("#%d", u.StringIndex)
			b.WriteString(s)
			b.WriteString(strings.Repeat(" ", l-len(s)))
			cm, err = u.ParseStringFromPool(m.CpInfo)
			if err != nil {
				return
			}
			cmShow = true
		case *IntegerInfo:
			b.WriteString(n)
			// todo
		case *Utf8Info:
			var s string
			if s, err = u2string(u.Bytes); err != nil {
				return
			}
			b.WriteString(s)
		case *NameAndType:
			s := fmt.Sprintf("#%d:#%d", u.NameIndex, u.DescriptorIndex)
			b.WriteString(s)
			b.WriteString(strings.Repeat(" ", l-len(s)))
			var name, desc string
			if name, desc, err = u.ParseFromPool(m.CpInfo); err != nil {
				return
			}
			cm = fmt.Sprintf("%s:%s", name, desc)
		}
		if cm != "" || cmShow {
			b.WriteString(fmt.Sprintf("// %s", cm))
		}
		b.WriteString("\n")
	}

	// flags
	fs := ParseClassAccessFlags(m.AccessFlags)
	b.WriteString("flags: ")
	fns := make([]string, 0)
	for _, fl := range fs {
		fns = append(fns, _cAccFm[fl])
	}
	b.WriteString(strings.Join(fns, ", "))
	b.WriteString("\n")

	// this class
	var qN string
	qN, err = qualifiedClassNameFromPool(m.CpInfo, m.ThisClass)
	if err != nil {
		return
	}
	b.WriteString(fmt.Sprintf("this class: %s\n", qN))

	// super class
	qN, err = qualifiedClassNameFromPool(m.CpInfo, m.SuperClass)
	if err != nil {
		return
	}
	b.WriteString(fmt.Sprintf("super class: %s\n", qN))

	// interfaces
	b.WriteString(fmt.Sprintf("interfaces count: %d\n", m.InterfacesCount))
	b.WriteString("interfaces: [")
	qns := make([]string, 0)
	for i := 0; i < int(m.InterfacesCount); i++ {
		qN, err = qualifiedClassName(m.CpInfo, m.Interfaces[i])
		qns = append(qns, qN)
	}
	b.WriteString(strings.Join(qns, ", "))
	b.WriteString("]\n")

	// fields
	b.WriteString(fmt.Sprintf("fields count: %d\n", m.FieldsCount))
	b.WriteString("fields: [\n")
	var fN string
	var fD string
	for i := 0; i < int(m.FieldsCount); i++ {
		f := m.Fields[i]
		fN, err = ui2string(m.CpInfo, f.NameIndex)
		if err != nil {
			return
		}
		fD, err = ui2string(m.CpInfo, f.DescriptorIndex)
		if err != nil {
			return
		}
		b.WriteString(fmt.Sprintf("	name: %s\n	desc: %s\n", fN, fD))

		// flags
		fs := ParseFieldAccessFlags(f.AccessFlags)
		fsNs := make([]string, 0)
		for _, fl := range fs {
			fsNs = append(fsNs, _fAccFm[fl])
		}
		b.WriteString(fmt.Sprintf("	flags: %s", strings.Join(fsNs, ",")))
		b.WriteString("\n")
		if i != (int(m.FieldsCount) - 1) {
			b.WriteString("\n")
		}
	}
	b.WriteString("]\n")

	// methods
	b.WriteString(fmt.Sprintf("methods count: %d\n", m.MethodsCount))
	b.WriteString("methods: [\n")
	for i := 0; i < int(m.MethodsCount); i++ {
		f := m.Methods[i]
		fN, err = ui2string(m.CpInfo, f.NameIndex)
		if err != nil {
			return
		}
		fD, err = ui2string(m.CpInfo, f.DescriptorIndex)
		if err != nil {
			return
		}
		b.WriteString(fmt.Sprintf("	name: %s\n	desc: %s\n", fN, fD))

		// flags
		fs := ParseMethodAccessFlags(f.AccessFlags)
		fsNs := make([]string, 0)
		for _, fl := range fs {
			fsNs = append(fsNs, _mAccFm[fl])
		}
		b.WriteString(fmt.Sprintf("	flags: %s", strings.Join(fsNs, ",")))
		b.WriteString("\n")
		if i != (int(m.MethodsCount) - 1) {
			b.WriteString("\n")
		}
	}
	b.WriteString("]\n")

	// attributes
	b.WriteString(fmt.Sprintf("attributes count: %d\n", m.AttributesCount))
	b.WriteString("attributes: [\n")
	for i := 0; i < int(m.AttributesCount); i++ {
		f := m.Attributes[i]
		fN, err = ui2string(m.CpInfo, f.AttributeNameIndex)
		if err != nil {
			return
		}
		switch fN {
		case _sourceFile:
			sfi, _ := u16(f.Info, 0)
			fD, err = ui2string(m.CpInfo, sfi)
			if err != nil {
				return
			}
		}
		b.WriteString(fmt.Sprintf("	name: %s\n	info: %s\n", fN, fD))
		if i != (int(m.AttributesCount) - 1) {
			b.WriteString("\n")
		}
	}
	b.WriteString("]\n")

	return b.String(), nil
}

func qualifiedClassNameFromPool(cp []ConstantInfo, i uint16) (name string, err error) {
	c, ok := cp[i].(*ClassInfo)
	if !ok {
		err = fmt.Errorf("index poiner to a non class info")
		return
	}
	return qualifiedClassName(cp, c)
}

func qualifiedClassName(cp []ConstantInfo, c *ClassInfo) (name string, err error) {
	name, err = c.ParseNameFromPool(cp)
	if err != nil {
		return
	}
	name = strings.Replace(name, "/", ".", -1)
	return
}

func u2string(b []byte) (res string, err error) {
	var rs []rune
	rs, err = DecodeRunes(b)
	if err != nil {
		return
	}
	res = string(rs)
	return
}

// FieldInfo field info.
type FieldInfo struct {
	AccessFlags     uint16
	NameIndex       uint16
	DescriptorIndex uint16
	AttributesCount uint16
	Attributes      []*AttributeInfo
}

func (m *FieldInfo) Read(b []byte, s int) (next int) {
	m.AccessFlags, next = u16(b, s)
	m.NameIndex, next = u16(b, next)
	m.DescriptorIndex, next = u16(b, next)
	m.AttributesCount, next = u16(b, next)
	m.Attributes = make([]*AttributeInfo, m.AttributesCount)
	for i := 0; i < int(m.AttributesCount); i++ {
		m.Attributes[i] = new(AttributeInfo)
		next = m.Attributes[i].Read(b, next)
	}
	return
}

// MethodInfo method info.
type MethodInfo struct {
	FieldInfo
}

// AttributeInfo attribute info.
type AttributeInfo struct {
	AttributeNameIndex uint16
	AttributeLength    uint32
	Info               []byte
}

func (m *AttributeInfo) Read(b []byte, s int) (next int) {
	m.AttributeNameIndex, next = u16(b, s)
	m.AttributeLength, next = u32(b, next)
	m.Info, next = bs(b, next, int(m.AttributeLength))
	return
}
