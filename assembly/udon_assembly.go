package assembly

import "fmt"

type UdonAssembly struct {
	ASM             string
	ProgramCounter  Addr
	IDCounter       int
	LabelDict       map[LabelName]Addr
	EventNames      []EventName
	ExportVars      []VarName
	VarTable        VarTable
	DefFuncTable    DefFuncTable
	UdonMethodTable UdonMethodTable
	EnvVars         []VarName
}

func (ua *UdonAssembly) AddInstComment(comment string) {
	return
}
func (ua *UdonAssembly) AddInst(bcodeSize Addr, inst string) {
	ua.ASM += fmt.Sprintf("%s\n", inst)
	ua.ProgramCounter = Addr(ua.ProgramCounter + bcodeSize)
	return
}
func (ua *UdonAssembly) MakeCodeSeg() string {
	ret_code := ".code_start\n\n"
	for _, eventName := range ua.EventNames {
		ret_code += fmt.Sprintf("    .export %s\n", eventName)
		ret_code += fmt.Sprintf("%s\n", ua.ASM)
		ret_code += fmt.Sprintf(".code_end\n")
	}
	return ret_code
}
func (ua *UdonAssembly) GetNextId(name string) string {
	ret_id := fmt.Sprintf("__%s_%d", name, ua.IDCounter)
	ua.IDCounter += 1
	return ret_id
}
func (ua *UdonAssembly) AddLabelCurrentAddr(label LabelName) {
	ua.LabelDict[label] = ua.ProgramCounter
	return
}
func (ua *UdonAssembly) Nop() {
	ua.AddInst(4, "NOP")
	return
}
func (ua *UdonAssembly) RemoveTop() {
	ua.AddInstComment("Remove Top")
	ua.AddInst(4, "POP")
	return
}
func (ua *UdonAssembly) PopVar(ret_value_name VarName) error {
	ua.AddInstComment(fmt.Sprintf("Pop(Push->Copy) %s", ret_value_name))
	ret_type := ua.VarTable.GetVarType(ret_value_name)
	if ret_type == nil {
		return fmt.Errorf("pop: Variable %s is not defined.", ret_value_name)
	}
	ua.PushVar(ret_value_name)
	ua.Copy()
	return nil
}
func (ua *UdonAssembly) PopVars(var_names []VarName) {
	return
}
func (ua *UdonAssembly) Push(addr Addr) {
	return
}
func (ua *UdonAssembly) PushVar(var_name VarName) {
	return
}
func (ua *UdonAssembly) PushVars(var_names []VarName) {
	return
}
func (ua *UdonAssembly) Copy() {
	return
}
func (ua *UdonAssembly) PushStr(_str string) {
	return
}
func (ua *UdonAssembly) Jump(addr Addr) {
	return
}
func (ua *UdonAssembly) JumpLabel(label LabelName) {
	return
}
func (ua *UdonAssembly) JumpIfFalse(addr Addr) {
	return
}
func (ua *UdonAssembly) JumpIfFalseLabel(label LabelName) {
	return
}
func (ua *UdonAssembly) JumpIndirect(var_name VarName) {
	return
}
func (ua *UdonAssembly) JumpRetAddr() {
	return
}
func (ua *UdonAssembly) Extern(extern_str ExternStr) {
	return
}
func (ua *UdonAssembly) End() {
	return
}
func (ua *UdonAssembly) CallExtern(extern_str ExternStr, arg_vars []VarName) {
	return
}
func (ua *UdonAssembly) Assign(dist_var_name VarName, src_var_name VarName) {
	return
}
func (ua *UdonAssembly) SetBool(var_name VarName, bool_num bool) {
	return
}
func (ua *UdonAssembly) SetUint32(var_name VarName, num int) {
	return
}
func (ua *UdonAssembly) GetAddr(label LabelName) Addr {
	return 0
}
func (ua *UdonAssembly) AddLabel(label LabelName, addr Addr) {
	return
}
func (ua *UdonAssembly) ReplaceTmpAdrr(code string) string {
	return ""
}
func (ua *UdonAssembly) CallDefFunc(func_name FuncName, arg_var_names []VarName) *VarName {
	return nil
}
func (ua *UdonAssembly) AddEvent(event_name EventName, def_arg_var_names []VarName, def_arg_types []UdonTypeName) {
	return
}
func (ua *UdonAssembly) EventHead(event_name EventName) {
	return
}
