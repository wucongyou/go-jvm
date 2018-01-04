package class

const (
	// class
	_cAccPublic    = 0x0001 // 0000000000000001
	_cAccFinal     = 0x0010 // 0000000000010000
	_cAccSuper     = 0x0020 // 0000000000100000
	_cAccInterface = 0x0200 // 0000001000000000
	_cAccAbstract  = 0x0400 // 0000010000000000
	_cAccSynthetic = 0x1000 // 0001000000000000

	// field
	_fAccPublic    = 0x0001 // 0000000000000001
	_fAccPrivate   = 0x0002 // 0000000000000010
	_fAccProtected = 0x0004 // 0000000000000100
	_fAccStatic    = 0x0008 // 0000000000001000
	_fAccFinal     = 0x0010 // 0000000000010000
	_fAccVolatile  = 0x0040 // 0000000001000000
	_fAccTransient = 0x0080 // 0000000010000000
	_fAccSynthetic = 0x1000 // 0001000000000000
	_fAccEnum      = 0x4000 // 0100000000000000

	// method
	_mAccPublic       = 0x0001 // 0000000000000001
	_mAccPrivate      = 0x0002 // 0000000000000010
	_mAccProtected    = 0x0004 // 0000000000000100
	_mAccStatic       = 0x0008 // 0000000000001000
	_mAccFinal        = 0x0010 // 0000000000010000
	_mAccSynchronized = 0x0020 // 0000000000100000
	_mAccBridge       = 0x0040 // 0000000001000000
	_mAccVarargs      = 0x0080 // 0000000010000000
	_mAccNative       = 0x0100 // 0000000100000000
	_mAccAbstract     = 0x0400 // 0000010000000000
	_mAccStrict       = 0x0800 // 0000100000000000
	_mAccSynthetic    = 0x1000 // 0001000000000000
)

var (
	// class
	_cAccFm = map[uint16]string{
		_cAccPublic:    "ACC_PUBLIC",
		_cAccFinal:     "ACC_FINAL",
		_cAccSuper:     "ACC_SUPER",
		_cAccInterface: "ACC_INTERFACE",
		_cAccAbstract:  "ACC_ABSTRACT",
		_cAccSynthetic: "ACC_SYNTHETIC",
	}

	// field
	_fAccFm = map[uint16]string{
		_fAccPublic:    "ACC_PUBLIC",
		_fAccProtected: "ACC_PROTECTED",
		_fAccStatic:    "ACC_STATIC",
		_fAccFinal:     "ACC_FINAL",
		_fAccVolatile:  "ACC_VOLATILE",
		_fAccTransient: "ACC_TRANSIENT",
		_fAccSynthetic: "ACC_SYNTHETIC",
		_fAccEnum:      "ACC_ENUM",
	}

	// method
	_mAccFm = map[uint16]string{
		_mAccPublic:       "ACC_PUBLIC",
		_mAccProtected:    "ACC_PROTECTED",
		_mAccStatic:       "ACC_STATIC",
		_mAccFinal:        "ACC_FINAL",
		_mAccSynchronized: "ACC_SYNCHRONIZED",
		_mAccBridge:       "ACC_BRIDGE",
		_mAccVarargs:      "ACC_VARARGS",
		_mAccNative:       "ACC_NATIVE",
		_mAccAbstract:     "ACC_ABSTRACT",
		_mAccSynthetic:    "ACC_SYNTHETIC",
	}
)

func ParseClassAccessFlags(f uint16) (fs []uint16) {
	fs = make([]uint16, 0)
	if f&_cAccPublic != 0 {
		fs = append(fs, _cAccPublic)
	}

	if f&_cAccFinal != 0 {
		fs = append(fs, _cAccFinal)
	}

	if f&_cAccSuper != 0 {
		fs = append(fs, _cAccSuper)
	}

	if f&_cAccInterface != 0 {
		fs = append(fs, _cAccInterface)
	}

	if f&_cAccAbstract != 0 {
		fs = append(fs, _cAccAbstract)
	}

	if f&_cAccSynthetic != 0 {
		fs = append(fs, _cAccSynthetic)
	}
	return
}

func ParseFieldAccessFlags(f uint16) (fs []uint16) {
	fs = make([]uint16, 0)
	if f&_fAccPublic != 0 {
		fs = append(fs, _fAccPublic)
	}

	if f&_fAccPrivate != 0 {
		fs = append(fs, _fAccPrivate)
	}

	if f&_fAccProtected != 0 {
		fs = append(fs, _fAccProtected)
	}

	if f&_fAccStatic != 0 {
		fs = append(fs, _fAccStatic)
	}

	if f&_fAccFinal != 0 {
		fs = append(fs, _fAccFinal)
	}

	if f&_fAccVolatile != 0 {
		fs = append(fs, _fAccVolatile)
	}

	if f&_fAccTransient != 0 {
		fs = append(fs, _fAccTransient)
	}

	if f&_fAccSynthetic != 0 {
		fs = append(fs, _fAccSynthetic)
	}

	if f&_fAccEnum != 0 {
		fs = append(fs, _fAccEnum)
	}
	return
}

func ParseMethodAccessFlags(f uint16) (fs []uint16) {
	fs = make([]uint16, 0)
	if f&_mAccPublic != 0 {
		fs = append(fs, _mAccPublic)
	}

	if f&_mAccPrivate != 0 {
		fs = append(fs, _mAccPrivate)
	}

	if f&_mAccProtected != 0 {
		fs = append(fs, _mAccProtected)
	}

	if f&_mAccStatic != 0 {
		fs = append(fs, _mAccStatic)
	}

	if f&_mAccFinal != 0 {
		fs = append(fs, _mAccFinal)
	}

	if f&_mAccSynchronized != 0 {
		fs = append(fs, _mAccSynchronized)
	}

	if f&_mAccBridge != 0 {
		fs = append(fs, _mAccBridge)
	}

	if f&_mAccVarargs != 0 {
		fs = append(fs, _mAccVarargs)
	}

	if f&_mAccNative != 0 {
		fs = append(fs, _mAccNative)
	}

	if f&_mAccAbstract != 0 {
		fs = append(fs, _mAccAbstract)
	}

	if f&_mAccStrict != 0 {
		fs = append(fs, _mAccStrict)
	}

	if f&_mAccSynthetic != 0 {
		fs = append(fs, _mAccSynthetic)
	}
	return
}
