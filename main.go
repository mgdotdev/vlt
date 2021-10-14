package main

import (
	"vlt/cli"
	"vlt/storage"
)

func main() {
	cmd, args := cli.Parse()
	key := *cli.Request("Please enter your vlt key:")

	userFile := storage.UserFile{Key: key}
	operation := storage.Operation{Cmd: cmd, Args: args, File: userFile}
	operation.Execute()
}
