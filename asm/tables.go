package asm

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type VarItem struct {
	VarName      VarName
	TypeName     UdonTypeName
	InitialValue string
}

// VarTable holds local, global variables and current function for contextual purposes
type VarTable struct {
	VarDict        []*VarItem
	GlobalVarNames []VarName
	CurrentFuncID  *LabelName
}

// Find varName in the vartable
func (vt *VarTable) Find(varName VarName) (*VarItem, bool) {
	for _, item := range vt.VarDict {
		if item.VarName == varName {
			return item, true
		}
	}
	return nil, false
}

// // UdonTypeItem holds the type name and initial value of a variable
// type UdonTypeItem struct {
// 	Name         UdonTypeName
// 	InitialValue string
// }

// NewVarTable returns a new var table for the assembler
func NewVarTable() *VarTable {
	t := &VarTable{
		VarDict:        []*VarItem{},
		GlobalVarNames: []VarName{},
		CurrentFuncID:  nil,
	}
	return t
}

// ResolveVarname returns the VarName for the current function, avoiding naming collisions across functions
func (vt *VarTable) ResolveVarname(varName VarName) (VarName, error) {
	tmpVarname := VarName(fmt.Sprintf("%s_%s", *vt.CurrentFuncID, varName))
	varItem, ok := vt.Find(varName)

	// ignore missing entry, create new var
	if vt.CurrentFuncID == nil && !ok {
		return "", errors.New("nil current func ID and var does not exist")
	}
	if vt.CurrentFuncID != nil && !ok {
		return tmpVarname, nil
	}
	return varItem.VarName, nil
}

// SetCurrentFuncID sets the current func ID for contextual execution
func (vt *VarTable) SetCurrentFuncID(label *LabelName) {
	if label == nil {
		vt.CurrentFuncID = nil
		return
	}
	vt.CurrentFuncID = label
}

// AddVarGlobal adds varName to the global variables
func (vt *VarTable) AddVarGlobal(varName VarName) error {
	for _, savedVarName := range vt.GlobalVarNames {
		if varName == savedVarName {
			return errors.New("global variable already registered")
		}
	}
	vt.GlobalVarNames = append(vt.GlobalVarNames, varName)
	return nil
}

// AddVar adds varName to the variable table
func (vt *VarTable) AddVar(varName VarName, typeName UdonTypeName, initValueStr string) error {
	_, ok := vt.Find(varName)
	if ok {
		return errors.New("variable already registered")
	}
	if typeName == UdonTypeName("Void") {
		return errors.New("variable is void type")
	}
	vt.VarDict = append(vt.VarDict, &VarItem{varName, typeName, initValueStr})

	return nil
}

// GetVarType returns the UdonType of the varName if it exists
func (vt *VarTable) GetVarType(varName VarName) (UdonTypeName, error) {
	_, ok := vt.Find(varName)
	if !ok {
		savedUdonType, ok := UdonTypes[varName]
		if !ok {
			return "", fmt.Errorf("variable not defined: %s", varName)
		}
		return savedUdonType, nil
	}
	item, ok := vt.Find(varName)
	if !ok {
		return "", fmt.Errorf("variable not defined: %s", varName)
	}
	return item.TypeName, nil
}

// ValidVarType returns true if varName has type assertVarType
// It checks if the variable exists
// It pulls the UdonType struct out of the variable table and returns the type
func (vt *VarTable) ValidVarType(varName VarName, assertVarType UdonTypeName) (bool, error) {
	if _, ok := vt.Find(varName); !ok {
		return false, fmt.Errorf("type not defined: %s", varName)
	}
	typeName, err := vt.GetVarType(varName)
	if err != nil {
		return false, fmt.Errorf("type not defined: %s", varName)
	}
	return assertVarType == typeName, nil

}

// MakeDataSeg generates the data block of the Udon Assembly
func (vt *VarTable) MakeDataSeg() (string, error) {
	dataStr := ".data_start\n\n"
	for _, varName := range vt.GlobalVarNames {
		if _, ok := vt.Find(varName); !ok {
			return "", errors.New("global var does not exist: " + string(varName))
		}
		dataStr += fmt.Sprintf("    .export %s\n", varName)
	}

	for _, v := range vt.VarDict {
		if v.TypeName == "VRCUdonCommonInterfacesIUdonEventReceiver" {
			dataStr += fmt.Sprintf("    %s: %%VRCUdonUdonBehaviour, %s\n", v.VarName, v.InitialValue)
		} else {
			dataStr += fmt.Sprintf("    %s: %%%s, %s\n", v.VarName, v.TypeName, v.InitialValue)
		}
	}
	dataStr += "\n.data_end\n\n"
	return dataStr, nil
}

// FnKey is used to fetch the saved func from the funcmap
type FnKey struct {
	FuncName FuncName
	ArgTypes string
}

// FnValue holds return types and var names from a funcmap
type FnValue struct {
	ReturnType UdonTypeName
	ArgNames   []VarName
}

// FuncMap is the funcmap, holding registered functions
type FuncMap map[FnKey]*FnValue

// Put will set or update the registered func in the funcmap
func (f FuncMap) Put(funcName FuncName, argTypes []UdonTypeName, retType UdonTypeName, argNames []VarName) {
	v := &FnValue{retType, argNames}
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	f[FnKey{funcName, strings.Join(argsStr, ",")}] = v
}

// Exists returns true if the func is registered
func (f FuncMap) Exists(funcName FuncName, argTypes []UdonTypeName) bool {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	_, ok := f[FnKey{funcName, strings.Join(argsStr, ",")}]
	return ok
}

// Get returns the registered funcName from the funcmap
func (f FuncMap) Get(funcName FuncName, argTypes []UdonTypeName) *FnValue {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	return f[FnKey{funcName, strings.Join(argsStr, ",")}]
}

// GetRetType returns the type of function with signature funcName and argTypes
func (f FuncMap) GetRetType(funcName FuncName, argTypes []UdonTypeName) (UdonTypeName, error) {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	exists := f.Exists(funcName, argTypes)
	if !exists {
		return "", fmt.Errorf("Function %s %s is not defined. Are the argument types correct?", funcName, strings.Join(argsStr, ","))
	}
	fn := f.Get(funcName, argTypes)
	return fn.ReturnType, nil
}

// GetFunctionID builds a string that represents the ID of that func given funcName and argTypes
func (f FuncMap) GetFunctionID(funcName FuncName, argTypes []UdonTypeName) LabelName {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	return LabelName(fmt.Sprintf(`%s__%s`, funcName, strings.Join(argsStr, "_")))
}

// MethodMap holds the methodmap of the assembler
type MethodMap map[MethodKey]*MethodValue

// MethodKey is used to build the key for the methodmap
type MethodKey struct {
	MethodKind UdonMethodKind
	ModuleName UdonTypeName
	MethodName UdonMethodName
	ArgTypes   string
}

// NewMethodKey returns a new method key
func NewMethodKey(
	methodKind UdonMethodKind,
	moduleName UdonTypeName,
	methodName UdonMethodName,
	argTypes []UdonTypeName,
) MethodKey {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	return MethodKey{
		methodKind,
		moduleName,
		methodName,
		strings.Join(argsStr, ","),
	}
}

// MethodValue maps to the external func
type MethodValue struct {
	TypeName  UdonTypeName
	ExternStr string
}

// Put adds and updates the method to the methodmap
func (umm MethodMap) Put(
	methodKind UdonMethodKind,
	moduleType UdonTypeName,
	methodName UdonMethodName,
	argTypes []UdonTypeName,
	value *MethodValue,
) {
	umm[NewMethodKey(
		methodKind,
		moduleType,
		methodName,
		argTypes,
	)] = value
}

// Exists returns true if key exists in the methodmap
func (umm MethodMap) Exists(
	methodKind UdonMethodKind,
	moduleType UdonTypeName,
	methodName UdonMethodName,
	argTypes []UdonTypeName,
) bool {
	_, ok := umm[NewMethodKey(
		methodKind,
		moduleType,
		methodName,
		argTypes,
	)]
	return ok
}

// Get returns the method stored in the methodmap
func (umm MethodMap) Get(
	methodKind UdonMethodKind,
	moduleType UdonTypeName,
	methodName UdonMethodName,
	argTypes []UdonTypeName,
) (*MethodValue, bool) {
	v, ok := umm[NewMethodKey(
		methodKind,
		moduleType,
		methodName,
		argTypes,
	)]
	return v, ok
}

// NewUdonMethodTable returns a methodmap that is prefilled with external methods
func NewUdonMethodTable(rdr io.Reader) (MethodMap, error) {
	methodMap, err := ParseExterns(rdr)
	if err != nil {
		return nil, fmt.Errorf("load method map: %w", err)
	}
	return methodMap, nil
}

// GetRetTypeExternStr returns the return type and external name of the method
func (umm MethodMap) GetRetTypeExternStr(
	methodKind UdonMethodKind,
	udonModuleType UdonTypeName,
	methodName UdonMethodName,
	argTypes []UdonTypeName) (*MethodValue, error) {
	exists := umm.Exists(
		methodKind,
		udonModuleType,
		methodName,
		argTypes,
	)
	if !exists {
		return nil, fmt.Errorf("method does not exist: %s", methodName)
	}
	method, ok := umm.Get(
		methodKind,
		udonModuleType,
		methodName,
		argTypes,
	)
	if !ok {
		return nil, errors.New("method not in map")
	}
	return method, nil
}
