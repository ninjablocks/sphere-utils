package main

import (
	"encoding/json"
	"fmt"
	"github.com/ninjasphere/go-ninja/config"
	"os"
)

func main() {
	flatten := false

	for _, a := range os.Args[1:] {
		if a == "--flatten" || a == "--flat" {
			flatten = true
		}
	}

	theMap := config.GetAll(flatten)
	if flatten {
		for k, v := range theMap {
			fmt.Printf("%s=%v\n", k, v)
		}
	} else {
		out, _ := json.MarshalIndent(theMap, "", "  ")
		fmt.Println(string(out))
	}
}
