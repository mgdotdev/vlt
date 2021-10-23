vlt
===

A CLI password vault written in Go
---

vlt is a command line tool for storing passwords in encrypted files. Files are stored on a per-user basis in $HOME/.local/lib/vlt and encrypted using a GCM block cypher. Simply build and run.

`vlt <cmd> <params>`

cmd is a string that matches any of the following:

### add
Aliases: ["add", "+", "-a"]

adds an entry to the database.

### build
Aliases: ["build", "-b", "mk"]

builds an entry w/ strong password

### edit
Aliases: ["edit", "-e"]

edits an entry from the database

### get
Aliases: ["get", "-g", "ls"]

gets an entry from the database

### remove
Aliases: ["remove", "rm", "-"]

removes an entry from the database

### load
Aliases: ["load", "-l", "<-"]

loads an unencrypted .csv file to the database

### dump
Aliases: ["dump", "-d", "->"]

dumps the database to an unencrypted csv file
unencrypted

see `vlt <cmd> --help` for details on each cmd.