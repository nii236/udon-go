package assembly

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type VarTable struct {
	VarDict        map[VarName]*UdonTypeItem
	GlobalVarNames []VarName
	CurrentFuncId  *string
}

type UdonTypeItem struct {
	Name         UdonTypeName
	InitialValue string
}

func NewVarTable() *VarTable {
	t := &VarTable{
		VarDict:        map[VarName]*UdonTypeItem{},
		GlobalVarNames: []VarName{},
		CurrentFuncId:  nil,
	}
	return t
}

func (vt *VarTable) ResolveVarname(varName VarName) VarName {
	tmpVarname := VarName(fmt.Sprintf("%s_%s", *vt.CurrentFuncId, varName))

	_, exists := vt.VarDict[tmpVarname]
	if vt.CurrentFuncId != nil && exists {
		return tmpVarname
	}
	return varName

}
func (vt *VarTable) AddVarGlobal(varName VarName) error {
	for _, savedVarName := range vt.GlobalVarNames {
		if varName == savedVarName {
			return errors.New("global variable already registered")
		}
	}
	vt.GlobalVarNames = append(vt.GlobalVarNames, varName)
	return nil
}
func (vt *VarTable) AddVar(varName VarName, typeName UdonTypeName, initValueStr string) error {
	exists := vt.ExistVar(varName)
	if exists {
		return errors.New("variable already registered")
	}
	if typeName == UdonTypeName("Void") {
		return errors.New("variable is void type")
	}
	vt.VarDict[varName] = &UdonTypeItem{typeName, initValueStr}
	return nil
}

func (vt *VarTable) GetVarType(varName VarName) (UdonTypeName, error) {
	exists := vt.ExistVar(varName)
	if !exists {
		savedUdonType, ok := UdonTypes[varName]
		if !ok {
			return "", fmt.Errorf("variable not defined: %s", varName)
		}
		return savedUdonType, nil
	}
	retType, ok := vt.VarDict[varName]
	if !ok {
		return "", fmt.Errorf("variable not defined: %s", varName)
	}
	return retType.Name, nil
}
func (vt *VarTable) ValidVarType(varName VarName, assertVarType UdonTypeName) bool {
	if !vt.ExistVar(varName) {
		fmt.Println("type not defined:", varName)
		return false
	}
	typeName, ok := vt.VarDict[varName]
	if !ok {
		fmt.Println("type not defined:", varName)
		return false
	}
	return assertVarType == typeName.Name

}
func (vt *VarTable) ExistVar(varName VarName) bool {
	_, ok := vt.VarDict[varName]
	return ok
}
func (vt *VarTable) MakeDataSeg() (string, error) {
	data_str := ".data_start\n\n"
	for _, varName := range vt.GlobalVarNames {
		if !vt.ExistVar(varName) {
			return "", errors.New("global var does not exist: " + string(varName))
		}
		data_str += fmt.Sprintf("    .export %s\n", varName)
	}

	for k, v := range vt.VarDict {
		if v.Name == "VRCUdonCommonInterfacesIUdonEventReceiver" {
			data_str += fmt.Sprintf("    %s: %%VRCUdonUdonBehaviour, %s\n", k, v.InitialValue)
		} else {
			data_str += fmt.Sprintf("    %s: %%%s, %s\n", k, v.Name, v.InitialValue)
		}
	}
	data_str += "\n.data_end\n\n"
	return data_str, nil
}
func (vt *VarTable) PrintDataSeg() {
	fmt.Println(vt.MakeDataSeg())
}

type FnKey struct {
	FuncName FuncName
	ArgTypes string
}
type FnValue struct {
	ReturnType UdonTypeName
	ArgNames   []VarName
}
type Funcs map[FnKey]*FnValue

func (f Funcs) Put(funcName FuncName, argTypes []UdonTypeName, value *FnValue) {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	f[FnKey{funcName, strings.Join(argsStr, ",")}] = value
}
func (f Funcs) Exists(funcName FuncName, argTypes []UdonTypeName) bool {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	_, ok := f[FnKey{funcName, strings.Join(argsStr, ",")}]
	return ok
}

func (f Funcs) Get(funcName FuncName, argTypes []UdonTypeName) *FnValue {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	return f[FnKey{funcName, strings.Join(argsStr, ",")}]
}

type DefFuncTable struct {
	FuncDict Funcs
}

func NewDefFuncTable() *DefFuncTable {
	return &DefFuncTable{Funcs{}}
}

func (ft *DefFuncTable) AddFunc(
	funcName FuncName,
	argTypes []UdonTypeName,
	retType UdonTypeName,
	argNames []VarName,
) {
	ft.FuncDict.Put(funcName, argTypes, &FnValue{retType, argNames})
}
func (ft *DefFuncTable) ExistFunc(funcName FuncName, argTypes []UdonTypeName) bool {
	return ft.FuncDict.Exists(funcName, argTypes)
}
func (ft *DefFuncTable) GetRetType(funcName FuncName, argTypes []UdonTypeName) (UdonTypeName, error) {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	exists := ft.FuncDict.Exists(funcName, argTypes)
	if !exists {
		return "", fmt.Errorf("Function %s %s is not defined. Are the argument types correct?", funcName, strings.Join(argsStr, ","))
	}
	fn := ft.FuncDict.Get(funcName, argTypes)
	return fn.ReturnType, nil
}
func (ft *DefFuncTable) GetFunctionId(funcName FuncName, argTypes []UdonTypeName) string {
	argsStr := []string{}
	for _, argType := range argTypes {
		argsStr = append(argsStr, string(argType))
	}
	return fmt.Sprintf(`%s__%s}`, funcName, strings.Join(argsStr, "_"))
}

type UdonMethodMap map[MethodKey]*MethodValue
type MethodKey struct {
	MethodKind UdonMethodKind
	ModuleName UdonTypeName
	MethodName UdonMethodName
	ArgTypes   string
}

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

type MethodValue struct {
	TypeName  UdonTypeName
	ExternStr string
}

func (umm UdonMethodMap) Put(
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
func (umm UdonMethodMap) Exists(
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
func (umm UdonMethodMap) Get(
	methodKind UdonMethodKind,
	moduleType UdonTypeName,
	methodName UdonMethodName,
	argTypes []UdonTypeName,
) *MethodValue {
	return umm[NewMethodKey(
		methodKind,
		moduleType,
		methodName,
		argTypes,
	)]
}

type UdonMethodTable struct {
	UdonMethodDict UdonMethodMap
}

func NewUdonMethodTable(rdr io.Reader) (*UdonMethodTable, error) {
	methodMap, err := ParseExterns(rdr)
	if err != nil {
		return nil, fmt.Errorf("load method map: %w", err)
	}
	t := &UdonMethodTable{
		UdonMethodDict: methodMap,
	}
	return t, nil
}

func (umt *UdonMethodTable) GetRetTypeExternStr(
	methodKind UdonMethodKind,
	udonModuleType UdonTypeName,
	methodName UdonMethodName,
	argTypes []UdonTypeName) (*MethodValue, error) {
	exists := umt.UdonMethodDict.Exists(
		methodKind,
		udonModuleType,
		methodName,
		argTypes,
	)
	if !exists {
		return nil, fmt.Errorf("method does not exist: %s", methodName)
	}
	method := umt.UdonMethodDict.Get(
		methodKind,
		udonModuleType,
		methodName,
		argTypes,
	)
	return method, nil
}
