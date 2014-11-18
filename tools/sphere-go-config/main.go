package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/ninjasphere/go-ninja/config"
)

var flatten = flag.Bool("flatten", false, "Flatten the config tree")

func main() {
	flag.Parse()

	out, _ := json.Marshal(config.GetAll(*flatten))

	fmt.Println(string(out))
}
