package cli

import (
	"flag"
	"os"
	"testing"
)

func _assert(t *testing.T) func(bool) {
	return func(item bool) {
		if item == false {
			t.Errorf("failed assertion")
		}
	}
}

func _reset_flags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestParseAdd(t *testing.T) {
	assert := _assert(t)
	_reset_flags()
	os.Args = []string{"vlt", "add", "-s", "google.com", "-u", "test_user", "-p", "test_password"}
	cmd, args := Parse()

	assert(cmd == "add")
	assert(args["source"] == "google.com")
	assert(args["username"] == "test_user")
	assert(args["password"] == "test_password")
}

func TestParseGet(t *testing.T) {
	assert := _assert(t)
	_reset_flags()

	os.Args = []string{"vlt", "get", "-s", "google.com", "-u", "test_user", "-p", "test_password"}
	cmd, args := Parse()

	assert(cmd == "get")
	assert(args["source"] == "google.com")
	assert(args["username"] == "test_user")
	assert(args["password"] == "test_password")
}

func TestParsePanic(t *testing.T) {
	assert := _assert(t)
	_reset_flags()

	catch := func() {
		r := recover()
		assert(r != nil)
	}
	defer catch()

	os.Args = []string{"vlt", "FAIL", "-s", "google.com", "-u", "test_user", "-p", "test_password"}
	Parse()
}
