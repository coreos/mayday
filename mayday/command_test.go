package mayday

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonexistentCommand(t *testing.T) {
	cmd := Command{Args: []string{"nonexistent"}}
	err := cmd.Run()
	assert.Equal(t, err.Error(), `could not find "nonexistent" in PATH`)
}

func TestCommand(t *testing.T) {
	cmd := Command{Args: []string{"df"}}
	err := cmd.Run()

	assert.Nil(t, err)
	assert.Contains(t, cmd.Content.String(), "Filesystem")
}

func TestCommandWithArgs(t *testing.T) {
	cmd := Command{Args: []string{"echo", "hello"}}

	err := cmd.Run()
	assert.Nil(t, err)

	assert.Equal(t, cmd.Content.String(), "hello\n")
}

func TestHeader(t *testing.T) {
	cmd := Command{Args: []string{"echo", "-e", "hello", "world", "testing"}}
	assert.Equal(t, cmd.outputName(), "echo_-e_hello_world_testing")

	cmd.Run() // cmd.Content needs to be populated for header()

	hdr := cmd.header()
	assert.Equal(t, hdr.Name, "mayday_commands/echo_-e_hello_world_testing")
}
