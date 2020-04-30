package asm

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type UdonAssembly struct {
	ASM            string
	ProgramCounter Addr
	IDCounter      int
	LabelDict      map[LabelName]Addr
	EventNames     []EventName
	ExportVars     []VarName
	VarTable       *VarTable
	FuncTable      FuncMap
	MethodTable    MethodMap
	EnvVars        []VarName
}

func NewUdonAssembly(rdr io.Reader) (*UdonAssembly, error) {
	umt, err := NewUdonMethodTable(rdr)
	if err != nil {
		return nil, fmt.Errorf("create udon method table: %w", err)
	}
	result := &UdonAssembly{
		ASM:            "",
		ProgramCounter: 0,
		IDCounter:      0,
		LabelDict:      map[LabelName]Addr{},
		EventNames:     []EventName{},
		ExportVars:     []VarName{},
		VarTable:       NewVarTable(),
		FuncTable:      FuncMap{},
		MethodTable:    umt,
		EnvVars:        []VarName{},
	}
	return result, nil
}

func (ua *UdonAssembly) AddInstComment(comment string) {
	return
}
func (ua *UdonAssembly) AddInst(bcodeSize Addr, inst string) {
	ua.ASM += fmt.Sprintf("        %s\n", inst)
	ua.ProgramCounter = Addr(ua.ProgramCounter + bcodeSize)
	return
}
func (ua *UdonAssembly) MakeCodeSeg() string {
	ret_code := ".code_start\n\n"
	for _, eventName := range ua.EventNames {
		ret_code += fmt.Sprintf("    .export %s\n", eventName)
	}
	ret_code += fmt.Sprintf("%s\n", ua.ASM)
	ret_code += fmt.Sprintf(".code_end\n")
	return ret_code
}
func (ua *UdonAssembly) GetNextId(name string) VarName {
	ret_id := fmt.Sprintf("__%s_%d", name, ua.IDCounter)
	ua.IDCounter += 1
	return VarName(ret_id)
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
	_, err := ua.VarTable.GetVarType(ret_value_name)
	if err != nil {
		return fmt.Errorf("PopVar: %w", err)
	}
	ua.PushVar(ret_value_name)
	ua.Copy()
	return nil
}
func (ua *UdonAssembly) PopVars(varNames []VarName) error {
	stringVarNames := []string{}
	for _, varName := range varNames {
		stringVarNames = append(stringVarNames, string(varName))
	}
	ua.AddInstComment(fmt.Sprintf("Pops %s", strings.Join(stringVarNames, ",")))
	for i := len(varNames) - 1; i > 0; i-- {
		varName := varNames[i]
		ua.PopVar(varName)
	}
	return nil
}
func (ua *UdonAssembly) Push(addr Addr) {
	ua.AddInst(Addr(8), fmt.Sprintf("PUSH, %x", addr))
	return
}
func (ua *UdonAssembly) PushVar(varName VarName) {
	ua.AddInst(Addr(8), fmt.Sprintf("PUSH, %s", varName))
	return
}
func (ua *UdonAssembly) PushVars(varNames []VarName) {
	stringVarNames := []string{}
	for _, varName := range varNames {
		stringVarNames = append(stringVarNames, string(varName))
	}
	ua.AddInstComment(fmt.Sprintf("pushes %s", strings.Join(stringVarNames, ",")))
	for _, varName := range varNames {
		ua.PushVar(varName)
	}
	return
}
func (ua *UdonAssembly) Copy() {
	ua.AddInst(Addr(4), "COPY")
	return
}
func (ua *UdonAssembly) PushStr(val string) {
	ua.AddInst(Addr(8), fmt.Sprintf(`PUSH, "%s"`, val))
	return
}
func (ua *UdonAssembly) Jump(addr Addr) {
	ua.AddInst(Addr(8), fmt.Sprintf("JUMP, %#x", addr))
	return
}
func (ua *UdonAssembly) JumpLabel(label LabelName) {
	ua.AddInst(Addr(8), fmt.Sprintf("JUMP, ###%s###", label))
	return
}
func (ua *UdonAssembly) JumpIfFalse(addr Addr) {
	ua.AddInst(Addr(8), fmt.Sprintf("JUMP_IF_FALSE, %#x", addr))
	return
}
func (ua *UdonAssembly) JumpIfFalseLabel(label LabelName) {
	ua.AddInst(Addr(8), fmt.Sprintf("JUMP_IF_FALSE, ###%s###", label))
	return
}
func (ua *UdonAssembly) JumpIndirect(varName VarName) {
	ua.AddInst(Addr(8), fmt.Sprintf("JUMP_INDIRECT, %s", varName))
	return
}
func (ua *UdonAssembly) JumpRetAddr() {
	ua.AddInst(Addr(8), fmt.Sprintf("JUMP_INDIRECT, ret_addr"))
	return
}
func (ua *UdonAssembly) Extern(extern_str ExternStr) {
	ua.AddInst(Addr(8), fmt.Sprintf(`EXTERN, "%s"`, extern_str))
	return
}
func (ua *UdonAssembly) End() {
	ua.AddInst(Addr(8), "JUMP, 0xFFFFFFFF")
	return
}
func (ua *UdonAssembly) CallExtern(extern_str ExternStr, argVars []VarName) {
	stringVarNames := []string{}
	for _, varName := range argVars {
		stringVarNames = append(stringVarNames, string(varName))
	}
	ua.AddInstComment(fmt.Sprintf("Call Extern %s[%s]", extern_str, strings.Join(stringVarNames, ",")))
	for _, arg := range argVars {
		ua.PushVar(arg)
	}
	ua.Extern(extern_str)
	return
}
func (ua *UdonAssembly) Assign(distVarName VarName, srcVarName VarName) error {
	// If the variable name on the right side is UdonTypeName,
	// just set the type of the variable on the left.
	existingVarType, srcVarNameexists := UdonTypes[srcVarName]
	if srcVarNameexists {
		ua.VarTable.AddVar(distVarName, existingVarType, "null")
		return nil
	}
	exists := ua.VarTable.ExistVar(distVarName)
	// fmt.Println(distVarName, srcVarName, exists)
	// If the left variable is undefined, define the variable.
	if !exists {
		srcVarType, err := ua.VarTable.GetVarType(srcVarName)
		if err != nil {
			return fmt.Errorf("get srcVarType: %w", err)
		}
		ua.AddInstComment(fmt.Sprintf("Declare %s", distVarName))
		ua.VarTable.AddVar(distVarName, srcVarType, "null")
	}
	ua.AddInstComment(fmt.Sprintf("%s = %s", distVarName, srcVarName))
	ua.PushVar(srcVarName)
	ua.PushVar(distVarName)
	ua.Copy()
	return nil
}
func (ua *UdonAssembly) SetBool(varName VarName, bool_num bool) {
	ua.PushStr(fmt.Sprintf("%v", bool_num))
	ua.PushVar(varName)
	ua.Extern(ExternStr("SystemBoolean.__Parse__SystemString__SystemBoolean"))
	return
}
func (ua *UdonAssembly) SetUint32(varName VarName, num int) {
	ua.AddInstComment(fmt.Sprintf("%s = %d", varName, num))
	constVarName := VarName(ua.GetNextId("const_uint32"))
	ua.VarTable.AddVar(constVarName, UdonTypeUInt32, strconv.Itoa(num))
	ua.PushVar(constVarName)
	ua.PushVar(varName)
	ua.Copy()
	return
}
func (ua *UdonAssembly) GetAddr(label LabelName) Addr {
	return ua.LabelDict[label]
}
func (ua *UdonAssembly) AddLabel(label LabelName, addr Addr) {
	ua.LabelDict[label] = addr
}
func (ua *UdonAssembly) ReplaceTmpAdrr(code string) string {
	result := []string{}
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if !strings.Contains(line, "###") {
			result = append(result, line)
			continue
		}
		r := regexp.MustCompile(".*###(.*)###.*")
		matches := r.FindStringSubmatch(line)
		if len(matches) < 2 {
			result = append(result, line)
			continue
		}
		// labelName := matches[1]
		// fmt.Println(line)
		line = r.ReplaceAllString(line, "hi")
		// line = strings.Replace(line, "###", "", -1)
		// addr := ua.GetAddr(LabelName(labelName))
		// line += strconv.Itoa(int(addr))
		result = append(result, line)

		// fmt.Sprintf("0x%d", addr))
	}
	return strings.Join(result, "\n")
}
func (ua *UdonAssembly) CallDefFunc(func_name FuncName, arg_var_names []VarName) (*VarName, error) {
	ua.AddInstComment(fmt.Sprintf("Call DefFunc %s%s", func_name, arg_var_names))
	arg_var_types := []UdonTypeName{}
	for _, argVarName := range arg_var_names {
		varType, err := ua.VarTable.GetVarType(argVarName)
		if err != nil {
			return nil, fmt.Errorf("CallDefFunc: %w", err)
		}
		arg_var_types = append(arg_var_types, varType)
	}
	retTypeName, err := ua.FuncTable.GetRetType(func_name, arg_var_types)
	if err != nil {
		return nil, fmt.Errorf("Get return type: %w", err)
	}
	retCallLabel := LabelName(ua.GetNextId("ret_call_label"))
	constRetAddr := VarName(ua.GetNextId("const_ret_addr"))
	retValue := VarName(ua.GetNextId("ret_value"))

	// Save current return address
	ua.PushVar(VarName("ret_addr"))
	// Save environment variables
	ua.PushVars(ua.EnvVars)
	// Save return address in order to return
	ua.VarTable.AddVar(
		VarName(constRetAddr),
		UdonTypeName("UInt32"),
		fmt.Sprintf("###%s###", retCallLabel),
	)
	// ua.Assign(VarName('ret_addr'), VarName(constRetAddr))
	ua.PushVar(VarName(constRetAddr))
	//Push arguments
	ua.PushVars(arg_var_names)
	//goto func label
	ua.JumpLabel(LabelName(ua.FuncTable.GetFunctionID(func_name, arg_var_types)))
	ua.AddLabelCurrentAddr(retCallLabel)
	if retTypeName != UdonTypeName("Void") {
		// pop ret_var_name
		ua.VarTable.AddVar(retValue, retTypeName, "null")
		ua.PopVar(retValue)
		// restore environment
		ua.PopVars(ua.EnvVars)
		// restore current return address
		ua.PopVar(VarName("ret_addr"))
		return &retValue, nil
	}
	// restore environment
	ua.PopVars(ua.EnvVars)
	// restore current return address
	ua.PopVar(VarName("ret_addr"))
	return nil, nil
}
func (ua *UdonAssembly) AddEvent(event_name EventName, def_arg_var_names []VarName, def_arg_types []UdonTypeName) error {
	savedEventItem, ok := EventTable[event_name]
	if !ok {
		// TODO: Add user event processing
		// (I still don't understand the specifications of user events)
		ua.EventNames = append(ua.EventNames, event_name)
		return nil
	}

	// Define the variables required for the event with arguments.
	if len(savedEventItem) != len(def_arg_var_names) {
		return fmt.Errorf("add_event: The required arguments for event %s and the number of defined arguments are different.", event_name)
	}

	for i, savedArgTuple := range savedEventItem {
		if savedArgTuple.UdonTypeName != def_arg_types[i] {
			return fmt.Errorf("add_event: The type of the argument of registered event %s is different from provided.", event_name)
		}
		ua.VarTable.AddVar(savedArgTuple.VarName, savedArgTuple.UdonTypeName, "null")
	}
	ua.EventNames = append(ua.EventNames, event_name)
	return nil
}
func (ua *UdonAssembly) EventHead(event_name EventName) {
	ua.ASM += fmt.Sprintf("    %s:\n", event_name)
}
