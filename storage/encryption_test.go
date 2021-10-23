package storage

import (
	"testing"
)

func _assert(t *testing.T) func(bool) {
	return func(item bool) {
		if item == false {
			t.Errorf("AssertionError")
		}
	}
}

func TestEncryption(t *testing.T) {
	assert := _assert(t)
	key := "encryption_key"
	expected := "fake_password"
	nonce := "asdfasdfasdf"
	temp := encrypt(expected, key, nonce)
	actual := decrypt(temp, key, nonce)
	assert(expected == actual)
}
