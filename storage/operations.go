package storage

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var headers []string = []string{"source", "username", "password"}

const numberStr string = "0123456789"
const punctuationStr string = "!#$%&()*+,-./:;<=>?@[]^_{|}~"
const lowerStr string = "abcdefghijklmnopqrstuvwxyz"
const upperStr string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
	case "edit":
		edit(opr.File, opr.Args)
	case "build":
		build(opr.File, opr.Args)
	case "remove":
		remove(opr.File, opr.Args)
	case "load":
		load(opr.File, opr.Args)
	case "dump":
		dump(opr.File, opr.Args)
	}
}

func add(file UserFile, args map[string]string) {
	entries := getEntries(file)

	if args["password"] == "_" {
		options := parseOptions(args["options"])
		delete(args, "options")
		args["password"] = makePassword(options)
	}

	entry := map[string]string{
		"source":   args["source"],
		"username": args["username"],
		"password": args["password"],
		"nonce":    AsNonce(uuid.NewString()),
	}

	entries.Add(entry)

	printEntry(len(entries.Items), entry)
	putEntries(file, entries)
}

func build(file UserFile, args map[string]string) {
	options := parseOptions(args["options"])
	delete(args, "options")
	args["password"] = makePassword(options)
	add(file, args)
}

func parseOptions(options string) map[string]string {
	byteOpt := []byte(options)

	optMap := make(map[string]string)
	re := regexp.MustCompile(`[^,]*=[^,]*`)
	opts := re.FindAll(byteOpt, -1)

	unpackSlice := func(items [][]byte) (string, string) {
		return string(items[0]), string(items[1])
	}

	for _, opt := range opts {
		key, val := unpackSlice(bytes.Split(opt, []byte("=")))
		optMap[key] = val
	}

	if optMap["m"] == "" {
		optMap["m"] = "alphanumeric!"
	}

	if optMap["l"] == "" {
		optMap["l"] = "20"
	}

	return optMap
}

func makePassword(options map[string]string) string {

	mode := options["m"]
	except := options["x"]
	length, err := strconv.Atoi(options["l"])
	check(err)

	var charPool string
	if mode == "uuid" {
		return uuid.NewString()
	}

	switch {
	case strings.Contains(mode, "alpha"):
		charPool += lowerStr + upperStr
	case strings.Contains(mode, "upper"):
		charPool += upperStr
	case strings.Contains(mode, "lower"):
		charPool += lowerStr
	}

	if strings.Contains(mode, "numeric") {
		charPool += numberStr
	}

	if strings.Contains(mode, "!") {
		charPool += punctuationStr
	}

	if len(charPool) == 0 {
		charPool += lowerStr + upperStr + numberStr + punctuationStr
	}

	exceptSlice := strings.Split(except, "")

	for _, item := range exceptSlice {
		charPool = strings.Replace(charPool, item, "", -1)
	}

	rand.Seed(time.Now().Unix())
	var result []byte

	for i := 0; i < length; i++ {
		result = append(result, charPool[rand.Intn(len(charPool))])
	}

	return string(result)
}

func edit(file UserFile, args map[string]string) {
	entries := getEntries(file)
	index := convIndex(args["index"], entries)
	delete(args, "index")

	if args["password"] == "_" {
		options := parseOptions(args["options"])
		delete(args, "options")
		args["password"] = makePassword(options)
	}

	entry := entries.Items[index]
	for _, key := range headers {
		item := args[key]
		if item != "" {
			entry[key] = item
		}
	}
	printEntry(index, entry)
	putEntries(file, entries)
}

func get(file UserFile, args map[string]string) {
	entries := getEntries(file)

	if args["index"] != "" {
		getByIndex(args["index"], entries)
		return
	}

	set := map[int]bool{}
	for i := 0; i < len(entries.Items); i++ {
		for _, key := range headers {
			if args[key] != "" && strings.Contains(
				entries.Items[i][key], args[key],
			) {
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
		printEntry(index, entry)
	}
}

func remove(file UserFile, args map[string]string) {
	entries := getEntries(file)
	index := convIndex(args["index"], entries)
	entries.Remove(index)
	putEntries(file, entries)
	fmt.Println("removed entry: " + strconv.Itoa(index))
}

func printEntry(index int, entry map[string]string) {
	str := []string{
		fmt.Sprint(index), ".\n",
		"  SOURCE: ", entry["source"], "\n",
		"USERNAME: ", entry["username"], "\n",
		"PASSWORD: ", entry["password"], "\n",
	}
	output := strings.Join(str, "")
	fmt.Println(output)
}

func getByIndex(index string, entries Entries) {
	i := convIndex(index, entries)
	printEntry(i, entries.Items[i])
}

func convIndex(index string, entries Entries) int {
	i, err := strconv.Atoi(index)
	check(err)

	if i < 0 {
		i = len(entries.Items) + i
	}
	return i
}

func load(file UserFile, args map[string]string) {
	data, err := os.ReadFile(args["path"])
	check(err)

	entries := Entries{}

	r := csv.NewReader(strings.NewReader(string(data)))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if isHeader(record) {
			continue
		}

		dict := make(map[string]string)
		dict["source"] = record[0]
		dict["username"] = record[1]
		dict["password"] = record[2]
		dict["nonce"] = AsNonce(uuid.NewString())

		entries.Add(dict)
	}

	putEntries(file, entries)
}

func isHeader(record []string) bool {
	if len(record) != len(headers) {
		return false
	}
	for i := 0; i < len(headers); i++ {
		if record[i] != headers[i] {
			return false
		}
	}
	return true
}

func dump(file UserFile, args map[string]string) {
	entries := getEntries(file)
	path := args["path"]
	fmt.Println(path)
	_, err := os.Stat(path)

	var f *os.File
	switch {
	case err == nil:
		confirm := *Request("a file of this name already exists. Continue? [y/n]")
		if confirm != "y" {
			fmt.Println("dump aborted")
			return
		}
		f, err = os.Create(path)
	case os.IsNotExist(err):
		f, err = os.Create(path)
	default:
		panic(err)
	}
	check(err)
	defer f.Close()

	f.WriteString(strings.Join(headers, ",") + "\n")
	for i := 0; i < len(entries.Items); i++ {
		entry := entries.Get(i)
		arr := []string{
			entry["source"],
			entry["username"],
			entry["password"],
		}
		line := strings.Join(arr, ",") + "\n"
		f.WriteString(line)
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
