package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/plenluno/jsm"
)

func main() {
	m := jsm.NewMachine()

	var p []jsm.Instruction
	j, _ := ioutil.ReadFile("../examples/fibonacci.json")
	json.Unmarshal(j, &p)

	res, _ := m.Run(p, []jsm.Value{jsm.NumberValue(35.0)})
	fmt.Println(res)
}
