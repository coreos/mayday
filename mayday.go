package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/coreos/mayday/mayday"

	"github.com/coreos/go-systemd/dbus"
)

const (
	dirPrefix = "mayday"
)

var (
	files = []mayday.File{
		{
			Path: "/proc/vmstat",
		},
		{
			Path: "/proc/meminfo",
			Link: "meminfo",
		},
		{
			Path: "/etc/os-release",
			Link: "os-release",
		},
		{
			Path: "/proc/mounts",
			Link: "mounts",
		},
	}
	commands = []mayday.Command{
		{
			Args: []string{"hostname"},
			Link: "hostname",
		},
		{
			Args: []string{"date"},
			Link: "date",
		},
		{
			Args: []string{"systemd-cgls"},
		},
		{
			Args: []string{"systemd-cgtop", "-n1"},
		},
		{
			Args: []string{"ps", "fauxwww"},
			Link: "ps",
		},
		{
			Args: []string{"lsmod"},
			Link: "lsmod",
		},
		{
			Args: []string{"lspci"},
			Link: "lspci",
		},
		{
			Args: []string{"lsof", "-b", "-M", "-n", "-l"},
			Link: "lsof",
		},
		{
			Args: []string{"blkid"},
		},
		{
			Args: []string{"btrfs", "fi", "show"},
		},
		{
			Args: []string{"df", "-al"},
			Link: "df",
		},
		{
			Args: []string{"ip", "-o", "addr", "show"},
			Link: "ip_addr_show",
		},
		{
			Args: []string{"ip", "-o", "link", "show"},
			Link: "ip_link_show",
		},
		{
			Args: []string{"ip", "-o", "route", "show"},
			Link: "ip_route_show",
		},
		{
			Args: []string{"netstat", "-neopa"},
			Link: "netstat",
		},
		{
			Args: []string{"df", "-ali"},
		},
		{
			Args: []string{"free", "-m"},
			Link: "free",
		},
		{
			Args: []string{"systemctl", "list-units", "-a"},
			Link: "all_units",
		},
		{
			Args: []string{"systemctl", "list-units", "--state=running"},
			Link: "running_units",
		},
		{
			Args: []string{"systemctl", "status", "etcd.service"},
			Link: "etcd_status",
		},
		{
			Args: []string{"systemctl", "status", "etcd2.service"},
			Link: "etcd2_status",
		},
		{
			Args: []string{"systemctl", "status", "fleet.service"},
			Link: "fleet_status",
		},
		{
			Args: []string{"systemctl", "status", "flanneld.service"},
			Link: "flanneld_status",
		},
	}
)

func getJournals(dir string) error {
	c, err := dbus.New()
	if err != nil {
		return err
	}

	defer c.Close()

	units, err := c.ListUnits()
	if err != nil {
		return err
	}

	pathre := regexp.MustCompile(`/usr/lib(32|64)?/systemd/system/.*\.service`)

	var svcs []string

	// build a list of units that live in /usr/lib/system/system
	for _, u := range units {
		if p, err := c.GetUnitProperty(u.Name, "FragmentPath"); err == nil {
			path := p.Value.Value().(string)

			if pathre.MatchString(path) {
				svcs = append(svcs, u.Name)
			}
		}
	}

	// get the journals of the units
	for _, svc := range svcs {
		daysago := 7

		logfile := path.Join(dir, fmt.Sprintf("%s.log", svc))

		log.Printf("collecting %d days of logs from %q", daysago, svc)

		out, err := os.Create(logfile)
		if err != nil {
			log.Printf("can't make log for %s: %s; skipping", svc, err)
			continue
		}

		cmd := exec.Command("journalctl", "--since", fmt.Sprintf("-%dd", daysago), "-l", "--utc", "--no-pager", "-u", svc)
		cmd.Stdout = out

		err = cmd.Run()
		out.Close()

		if err != nil {
			log.Printf("failed to dump log for %s: %s", svc, err)
			// journalctl seems to return a non-zero status when a service's log is empty,
			// so just delete the empty file.
			_ = os.Remove(logfile)
		}
	}

	return nil
}

func main() {
	now := time.Now().Format("200601021504-")
	ws, err := ioutil.TempDir("", dirPrefix+now)
	if err != nil {
		log.Fatalf("error creating output directory: %v", err)
	}
	log.Printf("Using workspace directory %q", ws)

	if err := os.MkdirAll(path.Join(ws, mayday.OutputDir), 0700); err != nil {
		log.Fatalf("error creating command output directory: %v", err)
	}

	journals := path.Join(ws, "journals")
	if err := os.MkdirAll(journals, 0700); err != nil {
		log.Fatalf("error creating command output directory: %v", err)
	}

	if err = getJournals(journals); err != nil {
		log.Printf("failed collecting journals: %s", err)
	}

	// TODO(jonboulle): parallelise
	// TODO(jonboulle): handle errors
	for _, f := range files {
		if err := f.Collect(ws); err != nil {
			fmt.Println(err)
		}
	}
	for _, c := range commands {
		if err := c.Run(ws); err != nil {
			fmt.Println(err)
		}
	}

	// Build the output tar file
	tfn := ws + ".tar.gz"
	f, err := os.OpenFile(tfn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("error opening file for output: %v", err)
	}
	gzw := gzip.NewWriter(f)
	if err := buildTarFile(gzw, ws); err != nil {
		log.Fatal(err.Error())
	}
	if err := gzw.Close(); err != nil {
		log.Fatalf("error closing zipfile: %v", err)
	}
	fmt.Printf("Output saved in %v\n", tfn)
	fmt.Println("All done!")
}

// buildTarFile writes a tar-formatted file containing the recursive contents
// of the given directory to the provided io.Writer
func buildTarFile(w io.Writer, dir string) error {
	tw := tar.NewWriter(w)
	wf := func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			// TODO(jonboulle): pass instead of failing?
			return err
		}
		var link string
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err = os.Readlink(p)
			link = strings.TrimPrefix(link, dir+"/")
		}
		if err != nil {
			return err
		}
		h, err := tar.FileInfoHeader(fi, link)
		if err != nil {
			return err
		}
		h.Name = strings.TrimPrefix(p, path.Dir(dir)+"/")
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		if fi.IsDir() || link != "" {
			return nil
		}
		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("could not open file for tar: %v", err)
		}
		if _, err := io.Copy(tw, f); err != nil {
			return fmt.Errorf("could not copy file: %v", err)
		}
		return nil
	}
	if err := filepath.Walk(dir, wf); err != nil {
		return fmt.Errorf("error collating output: %v", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("error closing tarfile: %v", err)
	}
	return nil
}
