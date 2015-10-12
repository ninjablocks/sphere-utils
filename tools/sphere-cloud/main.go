package main

import (
	"encoding/json"
	"fmt"
	"github.com/ninjasphere/go-ninja/cloud"
	"io/ioutil"
	"os"
)

const (
	USAGE = "sphere-cloud\n" +
		"   register email password name\n" +
		"   authenticate email password\n" +
		"   activtate token node\n" +
		"   get-tag token site tag\n" +
		"   post-tag token site tag < body\n" +
		"   put-tag token site tag < body\n"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, USAGE)
		os.Exit(1)
	}
	args = args[1:]

	switch args[0] {
	case "register":
		if len(args) < 4 {
			fmt.Fprintf(os.Stderr, USAGE)
			os.Exit(1)
		}
		if err := cloud.CloudAPI().RegisterUser(args[3], args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case "authenticate":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, USAGE)
			os.Exit(1)
		}
		if token, err := cloud.CloudAPI().AuthenticateUser(args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", token)
		}
	case "activate":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, USAGE)
			os.Exit(1)
		}
		if err := cloud.CloudAPI().ActivateSphere(args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case "get-tag":
		if len(args) < 4 {
			fmt.Fprintf(os.Stderr, USAGE)
			os.Exit(1)
		}
		var message json.RawMessage
		if err := cloud.CloudAPI().GetTag(args[1], args[2], args[3], &message); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		} else {
			if _, err := os.Stdout.Write(message); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
	case "post-tag", "put-tag":
		if len(args) < 4 {
			fmt.Fprintf(os.Stderr, USAGE)
			os.Exit(1)
		}
		if buffer, err := ioutil.ReadAll(os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		} else {
			data := map[string]interface{}{}
			if err := json.Unmarshal(buffer, &data); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			if err := cloud.CloudAPI().SetTag(args[1], args[2], args[3], data, args[0] == "put-tag"); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "unrecognized argument: %s\n", args[0])
		fmt.Fprintf(os.Stderr, USAGE)
		os.Exit(1)
	}
	os.Exit(0)
}
