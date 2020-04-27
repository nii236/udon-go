package assembly

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
