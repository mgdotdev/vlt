package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

func filename(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	name := hex.EncodeToString(h.Sum(nil))
	return name
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type UserFile struct {
	Key string
}

func (f *UserFile) Path() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".local", "lib", "vlt")

	name := filename(f.Key)
	path := filepath.Join(dir, name)
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

	text := encrypt(output.String(), f.Key)
	file.Write([]byte(text))
	file.Sync()
}

func (f *UserFile) Read() []byte {
	raw, err := os.ReadFile(f.Path())
	check(err)
	return []byte(decrypt(string(raw), f.Key))
}

func (f *UserFile) Write(output string) {
	err := os.WriteFile(f.Path(), []byte(encrypt(output, f.Key)), 0644)
	check(err)
}

type Entries struct {
	Items []map[string]string `json:"entries"`
}

func (e *Entries) Lock(key string) {
	for i := 0; i < len(e.Items); i++ {
		e.Items[i]["source"] = encrypt(e.Items[i]["source"], key)
		e.Items[i]["username"] = encrypt(e.Items[i]["username"], key)
		e.Items[i]["password"] = encrypt(e.Items[i]["password"], key)
	}
}

func (e *Entries) Unlock(key string) {
	for i := 0; i < len(e.Items); i++ {
		e.Items[i]["source"] = decrypt(e.Items[i]["source"], key)
		e.Items[i]["username"] = decrypt(e.Items[i]["username"], key)
		e.Items[i]["password"] = decrypt(e.Items[i]["password"], key)
	}
}
