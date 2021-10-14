package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Operation struct {
	Cmd  string
	Args map[string]string
	File UserFile
}

func (opr *Operation) Execute() {
	cmd := opr.Cmd

	switch cmd {
	case "add":
		add(opr.File, opr.Args)
	case "get":
		get(opr.File, opr.Args)
	}
}

func add(file UserFile, args map[string]string) {
	raw := file.Read()

	var entries Entries
	json.Unmarshal(raw, &entries)
	entries.Unlock(file.Key)

	item := map[string]string{
		"source":   args["source"],
		"username": args["username"],
		"password": args["password"],
	}

	entries.Items = append(entries.Items, item)
	entries.Lock(file.Key)

	res, err := json.Marshal(entries)
	check(err)

	var output bytes.Buffer
	err = json.Indent(&output, res, "", "\t")
	check(err)

	file.Write(output.String())
}

func get(file UserFile, args map[string]string) {
	raw := file.Read()

	var entries Entries
	json.Unmarshal(raw, &entries)
	entries.Unlock(file.Key)

	set := map[int]bool{}

	for i := 0; i < len(entries.Items); i++ {
		for key, value := range args {
			if value == entries.Items[i][key] {
				set[i] = true
			}
		}
	}

	if len(set) == 0 {
		fmt.Println("The query produced zero results.")
		return
	}

	for index := range set {
		entry := entries.Items[index]
		str := []string{
			fmt.Sprint(index), ".\n",
			"  SOURCE: ", entry["source"], "\n",
			"USERNAME: ", entry["username"], "\n",
			"PASSWORD: ", entry["password"], "\n",
		}

		output := strings.Join(str, "")
		fmt.Println(output)
	}
}
