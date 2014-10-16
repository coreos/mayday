package mayday

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	OutputDir = "mayday_commands"

	defaultTimeout = 30 * time.Second
)

// Command encapsulates a command (a list of arguments) to be run
type Command struct {
	Args []string // all of the arguments, e.g. ["free", "-m"]
	Link string   // short name to link to the output (optional), e.g. "free"
}

func (c *Command) outputFile() string {
	return strings.Join(c.Args, "_")
}

// Run runs the command, saving output to the given workspace
func (c *Command) Run(workspace string) error {
	// Sanitize provided arguments
	if len(c.Args) < 1 {
		return fmt.Errorf("cannot run empty Command")
	}
	name := c.Args[0]
	p, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("could not find %q in PATH", name)
	}

	// Set up the output file
	fn := path.Join(workspace, OutputDir, c.outputFile())
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error opening file for output: %v", err)
	}

	// Set up the actual Cmd to be run
	cmd := exec.Cmd{
		Path:   p,
		Args:   c.Args,
		Stdout: f,
		// TODO(jonboulle): something with stderr?
		// sosreport just appears to ignore it entirely.
	}

	// Launch the Cmd, and set up a timeout
	log.Printf("Running Command: %q\n", strings.Join(cmd.Args, " "))
	log.Printf("Saving output to %v\n", fn)
	cmd.Start()
	wc := make(chan error, 1)
	go func() {
		wc <- cmd.Wait()
	}()
	select {
	case <-time.After(defaultTimeout):
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Error killing Command: %v", err)
		}
		return fmt.Errorf("Timed out after %v running Command: %q", defaultTimeout, strings.Join(cmd.Args, " "))
	case err := <-wc:
		if err != nil {
			return err
		}
	}
	// If we get this far, the command succeeded. Huzzah!

	// If necessary, create a symlink
	return maybeCreateLink(c.Link, fn, workspace)
}
