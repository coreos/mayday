package mayday_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

var workspace string

func TestMain(m *testing.M) {
	// test setup
	// current_dir lives in /tmp/go-build
	current_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	workspace = current_dir + "/mayday_test"
	os.MkdirAll(workspace+"/mayday_commands", os.ModePerm)

	retCode := m.Run()

	os.Exit(retCode)
}
