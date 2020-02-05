package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

func prettyPrint(section string, data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("== %s ==\n", strings.ToTitle(section))
	fmt.Printf("%s \n", p)
}
