package asm_test

import (
	"fmt"
	"os"
	"testing"
	"udon-go/asm"
)

func TestVarTable(t *testing.T) {
	varTable := asm.NewVarTable()
	varTable.AddVar("aaa", "Int32", "100")
	varTable.AddVar("bbb", "Int32", "200")
	fmt.Println(varTable.MakeDataSeg())

	f, err := os.Open("./udon_funcs_data.txt")
	if err != nil {
		t.Errorf("open file: %s", err)
	}
	defer f.Close()
	udonMethodTable, err := asm.NewUdonMethodTable(f)
	if err != nil {
		t.Errorf("load udon method table: %v", err)
	}
	// pp.pprint(udonMethodTable.udon_method_dict)
	fmt.Println(udonMethodTable.GetRetTypeExternStr(
		"InstanceFunc",
		"ByteArray",
		"GetValue",
		[]asm.UdonTypeName{"Int32"},
	),
	)
	// ByteArray.GetValue Int32
}

func TestVarTable_AddVar(t *testing.T) {
	type fields struct {
		VarDict        map[asm.VarName]*asm.UdonTypeItem
		GlobalVarNames []asm.VarName
		CurrentFuncID  *asm.LabelName
	}
	type args struct {
		varName      asm.VarName
		typeName     asm.UdonTypeName
		initValueStr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"pristine", fields{map[asm.VarName]*asm.UdonTypeItem{}, []asm.VarName{}, nil}, args{"foo", asm.GoString, ""}, false},
		{"existing local var", fields{map[asm.VarName]*asm.UdonTypeItem{"foo": {asm.GoString, ""}}, []asm.VarName{}, nil}, args{"foo", asm.GoString, ""}, true},
		{"existing global var in var table", fields{map[asm.VarName]*asm.UdonTypeItem{"foo": {asm.GoString, ""}}, []asm.VarName{"foo"}, nil}, args{"foo", asm.GoString, ""}, true},
		{"existing global var not in var table", fields{map[asm.VarName]*asm.UdonTypeItem{}, []asm.VarName{"foo"}, nil}, args{"foo", asm.GoString, ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := &asm.VarTable{
				VarDict:        tt.fields.VarDict,
				GlobalVarNames: tt.fields.GlobalVarNames,
				CurrentFuncID:  tt.fields.CurrentFuncID,
			}
			if err := vt.AddVar(tt.args.varName, tt.args.typeName, tt.args.initValueStr); (err != nil) != tt.wantErr {
				t.Errorf("VarTable.AddVar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVarTable_GetVarType(t *testing.T) {
	type fields struct {
		VarDict        map[asm.VarName]*asm.UdonTypeItem
		GlobalVarNames []asm.VarName
		CurrentFuncID  *asm.LabelName
	}
	type args struct {
		varName asm.VarName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    asm.UdonTypeName
		wantErr bool
	}{
		{"pristine", fields{map[asm.VarName]*asm.UdonTypeItem{}, []asm.VarName{}, nil}, args{"foo"}, "", true},
		{"existing string", fields{map[asm.VarName]*asm.UdonTypeItem{"foo": {asm.GoString, ""}}, []asm.VarName{}, nil}, args{"foo"}, asm.GoString, false},
		{"existing int", fields{map[asm.VarName]*asm.UdonTypeItem{"foo": {asm.GoInt, ""}}, []asm.VarName{}, nil}, args{"foo"}, asm.GoInt, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := &asm.VarTable{
				VarDict:        tt.fields.VarDict,
				GlobalVarNames: tt.fields.GlobalVarNames,
				CurrentFuncID:  tt.fields.CurrentFuncID,
			}
			got, err := vt.GetVarType(tt.args.varName)
			if (err != nil) != tt.wantErr {
				t.Errorf("VarTable.GetVarType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VarTable.GetVarType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVarTable_AddVarGlobal(t *testing.T) {
	type fields struct {
		VarDict        map[asm.VarName]*asm.UdonTypeItem
		GlobalVarNames []asm.VarName
		CurrentFuncID  *asm.LabelName
	}
	type args struct {
		varName asm.VarName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"pristine", fields{map[asm.VarName]*asm.UdonTypeItem{}, []asm.VarName{}, nil}, args{"foo"}, false},
		{"existing global", fields{map[asm.VarName]*asm.UdonTypeItem{}, []asm.VarName{"foo"}, nil}, args{"foo"}, true},
		{"existing local", fields{map[asm.VarName]*asm.UdonTypeItem{"foo": {asm.GoString, ""}}, []asm.VarName{}, nil}, args{"foo"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := &asm.VarTable{
				VarDict:        tt.fields.VarDict,
				GlobalVarNames: tt.fields.GlobalVarNames,
				CurrentFuncID:  tt.fields.CurrentFuncID,
			}
			if err := vt.AddVarGlobal(tt.args.varName); (err != nil) != tt.wantErr {
				t.Errorf("VarTable.AddVarGlobal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVarTable_ValidVarType(t *testing.T) {
	type fields struct {
		VarDict        map[asm.VarName]*asm.UdonTypeItem
		GlobalVarNames []asm.VarName
		CurrentFuncID  *asm.LabelName
	}
	type args struct {
		varName       asm.VarName
		assertVarType asm.UdonTypeName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{"pristine", fields{map[asm.VarName]*asm.UdonTypeItem{}, []asm.VarName{}, nil}, args{"foo", asm.GoString}, false, true},
		{"existing wrong type", fields{map[asm.VarName]*asm.UdonTypeItem{"foo": {asm.GoInt, ""}}, []asm.VarName{"foo"}, nil}, args{"foo", asm.GoString}, false, false},
		{"existing matching type", fields{map[asm.VarName]*asm.UdonTypeItem{"foo": {asm.GoString, ""}}, []asm.VarName{}, nil}, args{"foo", asm.GoString}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := &asm.VarTable{
				VarDict:        tt.fields.VarDict,
				GlobalVarNames: tt.fields.GlobalVarNames,
				CurrentFuncID:  tt.fields.CurrentFuncID,
			}
			got, err := vt.ValidVarType(tt.args.varName, tt.args.assertVarType)
			if (err != nil) != tt.wantErr {
				t.Errorf("VarTable.ValidVarType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VarTable.ValidVarType() = %v, want %v", got, tt.want)
			}
		})
	}
}
