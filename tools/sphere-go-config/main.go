package main

import (
	"encoding/json"
	"fmt"

	"github.com/ninjasphere/go-ninja/config"
)

func main() {
	out, _ := json.Marshal(config.GetAll())

	fmt.Println(string(out))
}
