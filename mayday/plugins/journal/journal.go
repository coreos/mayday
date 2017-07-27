package journal

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/coreos/go-systemd/dbus"
	"github.com/coreos/mayday/mayday/tarable"
)

type SystemdJournal struct {
	name    string
	link    string        // currently never set to anything
	content *bytes.Buffer // the contents of the log, populated by Run()
}

type dbusStatus struct {
	unit     dbus.UnitStatus
	property *dbus.Property
}

var getJournals = func() ([]dbusStatus, error) {
	// return list of dbus Unit Statuses for use in ListJournals
	var units []dbus.UnitStatus
	var statuses []dbusStatus

	c, err := dbus.New()
	if err != nil {
		return statuses, err
	}

	defer c.Close()

	units, err = c.ListUnits()
	if err != nil {
		return statuses, err
	}
	for _, u := range units {
		if p, err := c.GetUnitProperty(u.Name, "FragmentPath"); err == nil {
			statuses = append(statuses, dbusStatus{unit: u, property: p})
		}
	}
	return statuses, nil
}

func List() ([]*SystemdJournal, error) {
	var svcs []*SystemdJournal

	statuses, err := getJournals()
	if err != nil {
		return svcs, err
	}

	pathre := regexp.MustCompile(`/usr/lib(32|64)?/systemd/system/.*\.service`)

	// build a list of units that live in /usr/lib/system/system
	for _, s := range statuses {
		path := s.property.Value.Value().(string)
		if pathre.MatchString(path) {
			svc := SystemdJournal{name: s.unit.Name}
			svcs = append(svcs, &svc)
		}
	}
	return svcs, nil
}

func (j *SystemdJournal) Content() *bytes.Buffer {
	if j.content == nil {
		j.Run()
	}
	return j.content
}

func (j *SystemdJournal) Name() string {
	return "/journals/" + j.name + ".log"
}

func (j *SystemdJournal) Header() *tar.Header {
	return tarable.Header(j.Content(), j.Name())
}

func (j *SystemdJournal) Link() string {
	return j.link
}

func (j *SystemdJournal) Run() error {
	var b bytes.Buffer
	j.content = &b
	writer := bufio.NewWriter(j.content)

	daysago := 7

	log.Printf("collecting %d days of logs from %q", daysago, j.name)

	cmd := exec.Command("journalctl", "--since", fmt.Sprintf("-%dd", daysago), "-l", "--utc", "--no-pager", "-u", j.name)
	cmd.Stdout = writer

	err := cmd.Run()

	if err != nil {
		log.Printf("failed to dump log for %s: %s", j.Name, err)
	}

	return err
}
