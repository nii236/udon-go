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

// GoNil is a convenience alias which maps a go type to a unity type
const GoNil UdonTypeName = "null"

// GoString is a convenience alias which maps a go type to a unity type
const GoString UdonTypeName = UdonTypeString

// GoBool is a convenience alias which maps a go type to a unity type
const GoBool UdonTypeName = UdonTypeBoolean

// GoInt16 is a convenience alias which maps a go type to a unity type
const GoInt16 UdonTypeName = UdonTypeInt16

// GoUint16 is a convenience alias which maps a go type to a unity type
const GoUint16 UdonTypeName = UdonTypeUInt16

// GoInt32 is a convenience alias which maps a go type to a unity type
const GoInt32 UdonTypeName = UdonTypeInt32

// GoUint32 is a convenience alias which maps a go type to a unity type
const GoUint32 UdonTypeName = UdonTypeUInt32

// GoInt64 is a convenience alias which maps a go type to a unity type
const GoInt64 UdonTypeName = UdonTypeInt64

// GoUint64 is a convenience alias which maps a go type to a unity type
const GoUint64 UdonTypeName = UdonTypeUInt64

// GoInt is a convenience alias which maps a go type to a unity type
const GoInt UdonTypeName = UdonTypeInt32

// GoUint is a convenience alias which maps a go type to a unity type
const GoUint UdonTypeName = UdonTypeUInt32

// GoFloat32 is a convenience alias which maps a go type to a unity type
const GoFloat32 UdonTypeName = UdonTypeDouble

// GoFloat64 is a convenience alias which maps a go type to a unity type
const GoFloat64 UdonTypeName = UdonTypeDouble

// GoRune is a convenience alias which maps a go type to a unity type
const GoRune UdonTypeName = UdonTypeChar

// GoByte is a convenience alias which maps a go type to a unity type
const GoByte UdonTypeName = UdonTypeByte
