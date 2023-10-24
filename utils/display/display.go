package display

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/alexeyco/simpletable"
)

type Details map[string]string

func PrintTable(obj interface{}, details Details) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	amorphous := map[string]interface{}{}
	err = json.Unmarshal(bytes, &amorphous)
	if err != nil {
		panic(err)
	}

	mainTable := simpletable.New()
	for row, field := range details {
		mainTable.Body.Cells = append(mainTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: row + ":"},
			{Text: fmt.Sprintf("%v", amorphous[field])},
		})
	}
	mainTable.SetStyle(simpletable.StyleCompact)

	return mainTable.String()
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func PrintList(input interface{}, details Details) string {
	table := simpletable.New()
	headers := []*simpletable.Cell{}
	for column := range details {
		headers = append(headers, &simpletable.Cell{Align: simpletable.AlignCenter, Text: strings.ToUpper(column)})
	}
	objs := interfaceSlice(input)
	for _, o := range objs {
		bytes, err := json.Marshal(o)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return ""
		}
		amorphous := map[string]interface{}{}
		err = json.Unmarshal(bytes, &amorphous)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return ""
		}

		row := []*simpletable.Cell{}
		for _, field := range details {
			row = append(row,
				&simpletable.Cell{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%v", amorphous[field])})
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}

	table.Header = &simpletable.Header{Cells: headers}
	table.SetStyle(simpletable.StyleCompact)
	return table.String()
}
