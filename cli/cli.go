package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func Parse() (string, map[string]string) {
	cmd := os.Args[1]
	switch cmd {
	case "add", "+", "-a":
		return "add", parseAdd()
	case "get", "-g":
		return "get", parseGet()
	default:
		panic(errors.New("unable to parse cli"))
	}
}

func Request(input string) *string {
	fmt.Println(input)
	fmt.Printf(">>> ")
	var item string
	fmt.Scan(&item)
	fmt.Printf("\n")
	return &item
}

func parseAdd() map[string]string {
	source := flag.String("s", "", "source of the entry")
	username := flag.String("u", "", "username of the entry")
	password := flag.String("p", "", "password of the entry")
	flag.CommandLine.Parse(os.Args[2:])

	dict := make(map[string]string)

	if *source == "" {
		source = Request("specify source:")
	}

	if *username == "" {
		username = Request("specify username:")
	}

	if *password == "" {
		password = Request("specify password:")
	}

	dict["source"] = *source
	dict["username"] = *username
	dict["password"] = *password
	return dict
}

func parseGet() map[string]string {
	source := flag.String("s", "", "source of the entry")
	username := flag.String("u", "", "username of the entry")
	password := flag.String("p", "", "password of the entry")

	flag.CommandLine.Parse(os.Args[2:])

	dict := make(map[string]string)
	dict["source"] = *source
	dict["username"] = *username
	dict["password"] = *password
	return dict
}
