package json_load

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func LoadJson(filename string, jsonobject interface{} ) {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(file, jsonobject)
}
