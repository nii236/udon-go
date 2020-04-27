package assembly_test

import (
	"fmt"
	"os"
	"testing"
	"udon-go/assembly"
)

func TestVarTable(t *testing.T) {
	varTable := assembly.NewVarTable()
	varTable.AddVar("aaa", "Int32", "100")
	varTable.AddVar("bbb", "Int32", "200")
	varTable.PrintDataSeg()

	f, err := os.Open("./udon_funcs_data.txt")
	if err != nil {
		t.Errorf("open file: %s", err)
	}
	defer f.Close()
	udonMethodTable, err := assembly.NewUdonMethodTable(f)
	if err != nil {
		t.Errorf("load udon method table: %v", err)
	}
	// pp.pprint(udonMethodTable.udon_method_dict)
	fmt.Println(udonMethodTable.GetRetTypeExternStr(
		"InstanceFunc",
		"ByteArray",
		"GetValue",
		[]assembly.UdonTypeName{"Int32"},
	),
	)
	// ByteArray.GetValue Int32
}
