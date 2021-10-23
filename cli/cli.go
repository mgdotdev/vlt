package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var optStr string = strings.Join(
	[]string{
		"build options for password generation:\n",
		"m=<string>\n\t",
		"sets the char select mode. can be\n\t",
		"either uuid, or a combination of\n\t",
		"[alpha, lower, upper, numeric, !].\n",
		"l=<number>\n\t",
		"length of the password.\n",
		"x=<string>\n\t",
		"characters to exclude from password.",
	},
	"",
)

func Parse() (string, map[string]string) {
	cmd := os.Args[1]
	switch cmd {
	case "add", "+", "-a":
		return "add", parseAdd()
	case "get", "-g", "ls":
		return "get", parseGet()
	case "edit", "-e":
		return "edit", parseEdit()
	case "build", "-b", "mk":
		return "build", parseBuild()
	case "remove", "rm", "-":
		return "remove", parseRemove()
	case "load", "-l", "<-":
		return "load", parseLoad()
	case "dump", "-d", "->":
		return "dump", parseDump()
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

func RequestKey(input string) string {
	fmt.Printf("Please enter your vlt key:\n>>> ")
	key, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Printf("\n\n")
	return string(key)
}

func parseAdd() map[string]string {
	source := flag.String("s", "", "source of the entry")
	username := flag.String("u", "", "username of the entry")
	password := flag.String("p", "", "password of the entry")
	options := flag.String("o", "", optStr)

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
	dict["options"] = *options
	return dict
}

func parseGet() map[string]string {
	source := flag.String("s", "", "source of the entry")
	username := flag.String("u", "", "username of the entry")
	password := flag.String("p", "", "password of the entry")
	index := flag.String("i", "", "integer of the entry")
	options := flag.String("o", "", optStr)

	flag.CommandLine.Parse(os.Args[2:])

	dict := make(map[string]string)
	dict["source"] = *source
	dict["username"] = *username
	dict["password"] = *password
	dict["index"] = *index
	dict["options"] = *options
	return dict
}

func parseEdit() map[string]string {
	dict := parseGet()
	if dict["password"] == "?" {
		dict["password"] = *Request("specify password:")
	}

	return dict
}

func parseBuild() map[string]string {
	source := flag.String("s", "", "source of the entry")
	username := flag.String("u", "", "username of the entry")
	options := flag.String("o", "", optStr)

	flag.CommandLine.Parse(os.Args[2:])

	if *source == "" {
		source = Request("specify source:")
	}

	if *username == "" {
		username = Request("specify username:")
	}

	dict := make(map[string]string)
	dict["source"] = *source
	dict["username"] = *username
	dict["options"] = *options
	return dict
}

func parseRemove() map[string]string {
	index := flag.String("i", "", "integer of the entry")

	flag.CommandLine.Parse(os.Args[2:])

	dict := make(map[string]string)
	dict["index"] = *index
	return dict
}

func parseLoad() map[string]string {
	path := os.Args[2]
	dict := make(map[string]string)
	dict["path"] = path
	return dict
}

func parseDump() map[string]string {
	return parseLoad()
}
