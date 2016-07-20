package command

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonexistentCommand(t *testing.T) {
	cmd := New([]string{"nonexistent"}, "")
	err := cmd.Run()
	assert.Equal(t, err.Error(), `could not find "nonexistent" in PATH`)
}

func TestCommand(t *testing.T) {
	cmd := New([]string{"df"}, "")
	err := cmd.Run()

	assert.Nil(t, err)

	content := new(bytes.Buffer)
	content.ReadFrom(cmd.Content())
	assert.Contains(t, content.String(), "Filesystem")
}

func TestCommandWithArgs(t *testing.T) {
	cmd := New([]string{"echo", "hello"}, "")

	err := cmd.Run()
	assert.Nil(t, err)

	content := new(bytes.Buffer)
	content.ReadFrom(cmd.Content())
	assert.Equal(t, content.String(), "hello\n")
}

func TestCommandHeader(t *testing.T) {
	cmd := New([]string{"echo", "-e", "hello", "world", "testing"}, "")
	assert.Equal(t, cmd.Name(), "/mayday_commands/echo_-e_hello_world_testing")

	cmd.Run() // cmd.Content needs to be populated for header()

	hdr := cmd.Header()
	assert.Equal(t, hdr.Name, "/mayday_commands/echo_-e_hello_world_testing")
}
