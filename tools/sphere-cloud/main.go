package main

import (
	"fmt"
	"github.com/ninjasphere/go-ninja/cloud"
	"os"
)

const (
	USAGE = "sphere-cloud\n" +
		"   register email password name\n" +
		"   authenticate email password\n" +
		"   activtate token node\n"
)

func main() {
	args := os.Args
	if len(args) < 1 {
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
	default:
		fmt.Fprintf(os.Stderr, "unrecognized argument: %s", args[0])
		fmt.Fprintf(os.Stderr, USAGE)
		os.Exit(1)
	}
	os.Exit(0)
}
