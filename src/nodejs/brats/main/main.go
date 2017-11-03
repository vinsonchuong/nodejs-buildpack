package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

func main() {
	dir := "/tmp/brats"
	nodejsVersion := "6.1.2"

	file, err := ioutil.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		panic(err)
	}
	obj := make(map[string]interface{})
	err = json.Unmarshal(file, &obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(obj["engines"])
	fmt.Println(reflect.TypeOf(obj["engines"]))
	engines, ok := obj["engines"].(map[string]interface{})
	fmt.Println(engines)
	if !ok {
		panic("conversion")
	}
	engines["node"] = nodejsVersion
	file, err = json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "package.json"), file, 0644)
	if err != nil {
		panic(err)
	}
}
