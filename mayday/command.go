package mayday

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

// Command encapsulates a command (a list of arguments) to be run
type Command struct {
	Args    []string      // all of the arguments, e.g. ["free", "-m"]
	Link    string        // short name to link to the output (optional), e.g. "free"
	Content *bytes.Buffer // the contents of the command, populated by Run()
}

func (c *Command) outputName() string {
	return strings.Join(c.Args, "_")
}

func (c *Command) header() *tar.Header {
	var header tar.Header
	header.Name = "mayday_commands/" + c.outputName()
	header.Size = int64(c.Content.Len())
	header.ModTime = time.Now()

	return &header

}

// Run runs the command, saving output to a Reader
func (c *Command) Run() error {

	var b bytes.Buffer
	c.Content = &b
	writer := bufio.NewWriter(c.Content)

	// Sanitize provided arguments
	if len(c.Args) < 1 {
		return fmt.Errorf("cannot run empty Command")
	}
	name := c.Args[0]
	p, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("could not find %q in PATH", name)
	}

	// Set up the actual Cmd to be run
	cmd := exec.Cmd{
		Path:   p,
		Args:   c.Args,
		Stdout: writer,
		// TODO(jonboulle): something with stderr?
		// sosreport just appears to ignore it entirely.
	}

	// Launch the Cmd, and set up a timeout
	log.Printf("Running command: %q\n", strings.Join(cmd.Args, " "))
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

	return nil
}
