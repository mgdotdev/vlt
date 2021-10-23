package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

func Filename(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	name := hex.EncodeToString(h.Sum(nil))
	return name
}

func AsNonce(str string) string {
	return str[:12]
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type UserFile struct {
	Key      string
	Filename string
}

func GetUserFile(key string) UserFile {
	name := Filename(key)
	return UserFile{Key: key, Filename: name}
}

func (f *UserFile) Nonce() string {
	return AsNonce(f.Filename)
}

func (f *UserFile) Path() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".local", "lib", "vlt")

	path := filepath.Join(dir, f.Filename)
	_, err := os.Stat(path)

	switch {
	case err == nil:
		return path
	case os.IsNotExist(err):
		f.build(dir, path)
		return path
	default:
		panic(err.Error())
	}
}

func (f *UserFile) build(dir string, path string) {
	os.MkdirAll(dir, 0755)
	file, err := os.Create(path)
	check(err)
	defer file.Close()

	entries := Entries{[]map[string]string{}}
	bytesJSON, _ := json.Marshal(entries)

	var output bytes.Buffer
	err = json.Indent(&output, bytesJSON, "", "\t")
	check(err)

	text := encrypt(output.String(), f.Key, f.Nonce())
	file.Write([]byte(text))
	file.Sync()
}

func (f *UserFile) Read() []byte {
	raw, err := os.ReadFile(f.Path())
	check(err)
	return []byte(decrypt(string(raw), f.Key, f.Nonce()))
}

func (f *UserFile) Write(output string) {
	err := os.WriteFile(f.Path(), []byte(encrypt(output, f.Key, f.Nonce())), 0644)
	check(err)
}

type Entries struct {
	Items []map[string]string `json:"entries"`
}

func (e *Entries) Add(entry map[string]string) {
	e.Items = append(e.Items, entry)
}

func (e *Entries) Remove(index int) {
	e.Items = append(e.Items[:index], e.Items[index+1:]...)
}

func (e *Entries) Get(index int) map[string]string {
	return e.Items[index]
}

func (e *Entries) Lock(key string) {
	for i := 0; i < len(e.Items); i++ {
		nonce := e.Items[i]["nonce"]
		e.Items[i]["source"] = encrypt(e.Items[i]["source"], key, nonce)
		e.Items[i]["username"] = encrypt(e.Items[i]["username"], key, nonce)
		e.Items[i]["password"] = encrypt(e.Items[i]["password"], key, nonce)
	}
}

func (e *Entries) Unlock(key string) {
	for i := 0; i < len(e.Items); i++ {
		nonce := e.Items[i]["nonce"]
		e.Items[i]["source"] = decrypt(e.Items[i]["source"], key, nonce)
		e.Items[i]["username"] = decrypt(e.Items[i]["username"], key, nonce)
		e.Items[i]["password"] = decrypt(e.Items[i]["password"], key, nonce)
	}
}

func getEntries(file UserFile) Entries {
	raw := file.Read()

	var entries Entries
	json.Unmarshal(raw, &entries)
	entries.Unlock(file.Key)
	return entries
}

func putEntries(file UserFile, entries Entries) {
	entries.Lock(file.Key)
	res, err := json.Marshal(entries)
	check(err)
	var output bytes.Buffer
	err = json.Indent(&output, res, "", "\t")
	check(err)
	file.Write(output.String())
}
