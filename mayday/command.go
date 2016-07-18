package mayday

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

// command encapsulates a command (a list of arguments) to be run
type Command struct {
	args    []string      // all of the arguments, e.g. ["free", "-m"]
	link    string        // short name to link to the output (optional), e.g. "free"
	content *bytes.Buffer // the contents of the command, populated by Run()
	Output  string        // name of command output file
}

func NewCommand(args []string, link string) *Command {
	c := new(Command)
	c.args = args
	c.link = link
	c.Output = "/mayday_commands/" + strings.Join(c.args, "_")
	return c
}

func (c *Command) Name() string {
	return c.Output
}

func (c *Command) Header() *tar.Header {
	// content needs to be populated before the header can be generated
	if c.content == nil {
		c.Run()
	}
	var header tar.Header
	header.Name = c.Name()
	header.Size = int64(c.content.Len())
	header.Mode = 0666
	header.ModTime = time.Now()

	return &header
}

func (c *Command) Content() io.Reader {
	if c.content == nil {
		c.Run()
	}
	return c.content
}

func (c *Command) Link() string {
	return c.link
}

// Run runs the command, saving output to a Reader
func (c *Command) Run() error {

	var b bytes.Buffer
	c.content = &b
	writer := bufio.NewWriter(c.content)

	// Sanitize provided arguments
	if len(c.args) < 1 {
		return fmt.Errorf("cannot run empty Command")
	}
	name := c.args[0]
	p, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("could not find %q in PATH", name)
	}

	// Set up the actual Cmd to be run
	cmd := exec.Cmd{
		Path:   p,
		Args:   c.args,
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
