package storage

import (
	"testing"
)

func TestFilename(t *testing.T) {
	assert := _assert(t)
	expected := Filename("test")
	actual := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	assert(expected == actual)
}

func TestEntries(t *testing.T) {
	assert := _assert(t)
	entries := new(Entries)

	dict := map[string]string{
		"source":   "source",
		"username": "username",
		"password": "password",
		"nonce":    "asdfasdfasdf",
	}

	entries.Items = append(entries.Items, dict)

	entries.Lock("test")
	entries.Unlock("test")

	assert(entries.Items[0]["source"] == dict["source"])
	assert(entries.Items[0]["username"] == dict["username"])
	assert(entries.Items[0]["password"] == dict["password"])
}
