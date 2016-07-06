package mayday

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"log"
	"os/exec"
	"regexp"
	"time"
)

type Journal struct {
	Name    string
	Content *bytes.Buffer // the contents of the log, populated by Get()
}

func GetJournals() ([]Journal, error) {
	var svcs []Journal

	c, err := dbus.New()
	if err != nil {
		return svcs, err
	}

	defer c.Close()

	units, err := c.ListUnits()
	if err != nil {
		return svcs, err
	}

	pathre := regexp.MustCompile(`/usr/lib(32|64)?/systemd/system/.*\.service`)

	// build a list of units that live in /usr/lib/system/system
	for _, u := range units {
		if p, err := c.GetUnitProperty(u.Name, "FragmentPath"); err == nil {
			path := p.Value.Value().(string)

			if pathre.MatchString(path) {
				svcs = append(svcs, Journal{Name: u.Name})
			}
		}
	}
	return svcs, nil
}

func (j *Journal) header() *tar.Header {
	var header tar.Header
	header.Name = "/journals/" + j.Name + ".log"
	header.Size = int64(j.Content.Len())
	header.ModTime = time.Now()

	return &header
}

func (j *Journal) Get() error {

	var b bytes.Buffer
	j.Content = &b
	writer := bufio.NewWriter(j.Content)

	daysago := 7

	log.Printf("collecting %d days of logs from %q", daysago, j.Name)

	cmd := exec.Command("journalctl", "--since", fmt.Sprintf("-%dd", daysago), "-l", "--utc", "--no-pager", "-u", j.Name)
	cmd.Stdout = writer

	err := cmd.Run()

	if err != nil {
		log.Printf("failed to dump log for %s: %s", j.Name, err)
	}

	return err
}
