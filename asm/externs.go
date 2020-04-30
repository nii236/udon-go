package asm

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func ParseExterns(rdr io.Reader) (MethodMap, error) {

	scan := bufio.NewScanner(rdr)

	result := MethodMap{}
	bad := 0
	for scan.Scan() {
		txt := scan.Text()
		// fmt.Println("process:", txt)
		key, value, err := parseRecord(txt)
		if err != nil {
			bad++
			continue
		}
		result[*key] = value
	}
	// fmt.Println("Processed records:", len(result))
	// fmt.Println("Bad records:", bad)
	return result, nil
}

func parseRecord(in string) (*MethodKey, *MethodValue, error) {
	r := regexp.MustCompile(`^\('(.*?)', '(.*?)', '(.*?)', \((.*?)\).*: \('(.*?)', '(.*?)'`)
	result := r.FindStringSubmatch(in)
	if len(result) < 5 {
		return nil, nil, errors.New("Bad length: " + strconv.Itoa(len(result)))
	}
	argTypes := result[4]
	argTypes = strings.ReplaceAll(argTypes, " ", "")
	argTypes = strings.ReplaceAll(argTypes, "'", "")
	numCommas := strings.Count(argTypes, ",")
	if numCommas <= 1 {
		argTypes = strings.ReplaceAll(argTypes, ",", "")
	}
	return &MethodKey{
			MethodKind: UdonMethodKind(result[1]),
			ModuleName: UdonTypeName(result[2]),
			MethodName: UdonMethodName(result[3]),
			ArgTypes:   argTypes,
		}, &MethodValue{
			TypeName:  UdonTypeName(result[5]),
			ExternStr: result[6],
		}, nil
}
