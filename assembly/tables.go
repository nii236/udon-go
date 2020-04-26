package assembly

import (
	"errors"
	"fmt"
)

type VarTable struct {
	VarDict        map[VarName]*UdonTypeName
	GlobalVarNames []VarName
	CurrentFuncId  *string
}

func NewVarTable() *VarTable {
	t := &VarTable{
		VarDict:        map[VarName]*UdonTypeName{},
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
func (vt *VarTable) AddVar(varName VarName, typeName *UdonTypeName, initValueStr string) {
	return
}

func (vt *VarTable) GetVarType(varName VarName) *UdonTypeName {
	return &UdonTypeName{}
}
func (vt *VarTable) ValidVarType(varName VarName, assertVarType *UdonTypeName) ValidStatus {
	return ""
}
func (vt *VarTable) ExistVar(varName VarName) bool {
	return false
}
func (vt *VarTable) MakeDataSeg() (string, error) {
	data_str := ".data_start\n\n"
	for _, varName := range vt.GlobalVarNames {
		if !vt.ExistVar(varName) {
			return "", errors.New("global var does not exist")
		}
		data_str += fmt.Sprintf(".export %s\n", varName)
	}

	for k, v := range vt.VarDict {
		if v.Name == "VRCUdonCommonInterfacesIUdonEventReceiver" {
			data_str += fmt.Sprintf("%s: %%VRCUdonUdonBehaviour, %s\n", k, v.InitValue)
		} else {
			data_str += fmt.Sprintf("%s: %%%s, %s\n", k, v.Name, v.InitValue)
		}
	}
	data_str += "\n.data_end\n\n"
	return data_str, nil
}
func (vt *VarTable) PrintDataSeg() {
	return
}

type DefFuncTable struct {
}

func (ft *DefFuncTable) AddFunc(funcName FuncName, argTypes []*UdonTypeName) {
	return
}
func (ft *DefFuncTable) ExistFunc(funcName FuncName, argTypes []*UdonTypeName) bool {
	return false
}
func (ft *DefFuncTable) GetRetType(funcName FuncName, argTypes []*UdonTypeName) *UdonTypeName {
	return nil
}
func (ft *DefFuncTable) GetFunctionId(funcName FuncName, argTypes []*UdonTypeName) string {
	return ""
}

type UdonMethod struct {
	UdonMethodKind UdonMethodKind
	UdonModuleName UdonModuleName
	UdonMethodName UdonMethodName
	UdonTypeName   []*UdonTypeName
}
type UdonMethodReturn struct {
	UdonTypeName *UdonTypeName
	ExternStr    ExternStr
}
type UdonMethodTable struct {
	UdonMethodDict map[*UdonMethod]UdonMethodReturn
}

func (umt *UdonMethodTable) GetRetTypeExternStr(
	methodKind UdonMethodKind,
	udonModuleType *UdonTypeName,
	methodName UdonMethodName,
	argTypes []*UdonTypeName) (*UdonTypeName, ExternStr) {
	return nil, ""
}
