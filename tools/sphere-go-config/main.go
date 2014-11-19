package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/ninjasphere/go-ninja/config"
)

var flatten = flag.Bool("flatten", false, "Flatten the config tree")
var flat = flag.Bool("flat", false, "Flatten the config tree. Alt")

func main() {
	flag.Parse()

	out, _ := json.Marshal(config.GetAll(*flatten || *flat))

	fmt.Println(string(out))
}
