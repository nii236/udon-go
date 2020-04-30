package asm

type VarName string
type FuncName string
type LabelName string
type EventName string
type UdonModuleName string
type UdonMethodName string
type ExternStr string
type UdonTypeName string

type Addr int

type ValidStatus string

const NOT_EXIST ValidStatus = "NOT_EXIST"
const VALID ValidStatus = "VALID"
const NOT_VALID ValidStatus = "NOT_VALID"

type UdonMethodKind string

const STATIC_FUNC UdonMethodKind = "StaticFunc"
const INSTANCE_FUNC UdonMethodKind = "InstanceFunc"
const CONSTRUCTOR UdonMethodKind = "Constructor"
const UNKNOWN UdonMethodKind = "Unknown"

// GoString is a convenience alias which maps a go type to a unity type
const GoString = UdonTypeString

// GoBool is a convenience alias which maps a go type to a unity type
const GoBool = UdonTypeBoolean

// GoInt16 is a convenience alias which maps a go type to a unity type
const GoInt16 = UdonTypeInt16

// GoUint16 is a convenience alias which maps a go type to a unity type
const GoUint16 = UdonTypeUInt16

// GoInt32 is a convenience alias which maps a go type to a unity type
const GoInt32 = UdonTypeInt32

// GoUint32 is a convenience alias which maps a go type to a unity type
const GoUint32 = UdonTypeUInt32

// GoInt64 is a convenience alias which maps a go type to a unity type
const GoInt64 = UdonTypeInt64

// GoUint64 is a convenience alias which maps a go type to a unity type
const GoUint64 = UdonTypeUInt64

// GoInt is a convenience alias which maps a go type to a unity type
const GoInt = UdonTypeInt32

// GoUint is a convenience alias which maps a go type to a unity type
const GoUint = UdonTypeUInt32

// GoFloat32 is a convenience alias which maps a go type to a unity type
const GoFloat32 = UdonTypeDouble

// GoFloat64 is a convenience alias which maps a go type to a unity type
const GoFloat64 = UdonTypeDouble

// GoRune is a convenience alias which maps a go type to a unity type
const GoRune = UdonTypeChar

// GoByte is a convenience alias which maps a go type to a unity type
const GoByte = UdonTypeByte
