package storage

import (
	"testing"
)

func _assert(t *testing.T) func(bool) {
	return func(item bool) {
		if item == false {
			t.Errorf("failed assertion")
		}
	}
}

func TestEncryption(t *testing.T) {
	assert := _assert(t)
	key := "encryption_key"
	expected := "fake_password"
	temp := encrypt(expected, key)
	actual := decrypt(temp, key)
	assert(expected == actual)
}
